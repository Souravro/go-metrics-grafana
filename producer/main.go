package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"producer/helper"
	"producer/producer_structs"
	"strings"
	"sync"
	"time"

	"github.com/Shopify/sarama"
)

var (
	brokers           = "broker_1:9092"
	topic             = "user_details_1"
	producerConfig    producer_structs.ProducerConfig
	productionCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "producer",
		Name:      "message_produced",
		Help:      "Counter for message produced",
	}, []string{"id"})
)

const (
	ProducerConfigFilename = "config.json"
)

func registerPrometheusMetrics() {
	prometheus.MustRegister(productionCounter)
}

func createConfig() *sarama.Config {
	config := sarama.NewConfig()
	config.Producer.Return.Errors = true
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll

	return config
}

func main() {
	producerConfig = helper.LoadProducerConfiguration(os.Getenv("APP_HOME") + "/config/" + ProducerConfigFilename)
	log.Printf("Producer Config: [%v]", producerConfig)

	log.Println("Starting a new Sarama producer...")
	ctx, cancel := context.WithCancel(context.Background())

	// Register Prometheus custom metrics
	http.Handle("/metrics", promhttp.Handler())
	registerPrometheusMetrics()

	// Start http server
	fmt.Printf("Starting server at port 8181...\n")
	go func() {
		if err := http.ListenAndServe(":8181", nil); err != nil {
			log.Fatal(err)
		}
	}()

	config := createConfig()
	producer, err := sarama.NewSyncProducer(strings.Split(brokers, ","), config)
	if err != nil {
		log.Printf("Producer. Error in creating producer. Error: [%v]", err)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			default:
				produceRecord(producer)
				time.Sleep(time.Duration(producerConfig.MessageInterval) * time.Millisecond)
			}
			//break
		}
	}()

	wg.Wait()
	cancel()
}

func getEncodedMessage() []byte {
	// Randomly create a json object of type Message, convert and return in []byte
	rand.Seed(time.Now().UnixNano())
	return encodeMessage(producer_structs.Message{
		Id:    producerConfig.UniqueIds[rand.Intn(len(producerConfig.UniqueIds))],
		Value: math.Round(producerConfig.ValuesMin+rand.Float64()*(producerConfig.ValuesMax-producerConfig.ValuesMin)*100) / 100,
	})
	//return encodeMessage(structs.Message{
	//	Id:    "1330",
	//	Value: 100.5,
	//})
}

func encodeMessage(msg interface{}) []byte {
	val, er := json.Marshal(msg)
	if er != nil {
		log.Printf("Error in converting struct to json. Error: [%v]", val)
	}

	return val
}

func getDecodedMessage(msg []byte) (producer_structs.Message, error) {
	var value producer_structs.Message
	if err := json.Unmarshal(msg, &value); err != nil {
		return producer_structs.Message{}, err
	}
	return value, nil
}

func produceRecord(producer sarama.SyncProducer) {
	// Produce records
	msgBytes := getEncodedMessage()
	producerMsg := &sarama.ProducerMessage{Topic: topic, Key: nil, Value: sarama.StringEncoder(msgBytes)}
	partition, offset, er := producer.SendMessage(producerMsg)
	if er != nil {
		log.Printf("Producer. Unable to Send Message to topic. Error: [%v]", er)
		return
	}
	log.Printf("Producer: message successfully published: produced message- [%v]. Partition: [%v]. Offset: [%v]", string(msgBytes), partition, offset)

	// Update production counter metric
	decodedMessage, _ := getDecodedMessage(msgBytes)
	productionCounter.WithLabelValues(decodedMessage.Id).Inc()
}
