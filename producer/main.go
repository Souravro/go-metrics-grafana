package main

import (
	"context"
	"encoding/json"
	"log"
	"math"
	"math/rand"
	"os"
	"producer/producer/producer_structs"
	"producer/structs"
	"strings"
	"sync"
	"time"

	"github.com/Shopify/sarama"
)

type Message struct {
	Id    string  `json:"id"`
	Value float64 `json:"value"`
}

var (
	brokers        = "0.0.0.0:8097"
	topic          = "user_details_1"
	producerConfig producer_structs.ProducerConfig
	commonConfig   structs.CommonConfig
)

func createConfig() *sarama.Config {
	config := sarama.NewConfig()
	config.Producer.Return.Errors = true
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll

	return config
}

func loadCommonConfiguration(file string) (common structs.CommonConfig) {
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		log.Printf("Error in opening common config file. Error: [%v]", err)
		return
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&common)

	return common
}

func loadProducerConfiguration(file string) (common producer_structs.ProducerConfig) {
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		log.Printf("Error in opening producer config file. Error: [%v]", err)
		return
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&common)

	return common
}

func main() {
	// load common config and producer config in memory
	commonConfig = loadCommonConfiguration("config/common.json")
	log.Printf("Common config: [%v]", commonConfig)

	producerConfig = loadProducerConfiguration("producer/config/config.json")
	log.Printf("Producer Config: [%v]", producerConfig)

	log.Println("Starting a new Sarama producer...")
	ctx, cancel := context.WithCancel(context.Background())

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
		}
	}()

	wg.Wait()
	cancel()
}

func getEncodedMessage() []byte {
	// Randomly create a json object of type Message, convert and return in []byte
	rand.Seed(time.Now().Unix())
	return encodeMessage(Message{
		Id:    commonConfig.UniqueIds[rand.Intn(len(commonConfig.UniqueIds))],
		Value: math.Round(commonConfig.ValuesMin + rand.Float64()*(commonConfig.ValuesMax-commonConfig.ValuesMin)),
	})
}

func encodeMessage(msg interface{}) []byte {
	val, er := json.Marshal(msg)
	if er != nil {
		log.Printf("Error in converting struct to json. Error: [%v]", val)
	}

	return val
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
}
