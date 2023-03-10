package main

import (
	"consumer/consumer_structs"
	"consumer/handler"
	"consumer/routes"
	"consumer/store"
	"context"
	"encoding/json"
	"fmt"
	"github.com/dgraph-io/badger"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/Shopify/sarama"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Consumer represents a Sarama consumer group consumer
type Consumer struct {
	ready chan bool
}

// Sarama configuration options
var (
	brokers            = "broker_1:9092"
	topic              = "user_details_1"
	group              = "user_group_1"
	oldest             = true
	storageSvc         *store.StorageService
	consumptionCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "consumer",
		Name:      "message_consumed",
		Help:      "Counter for message consumed",
	}, []string{"id"})
)

func registerPrometheusMetrics() {
	prometheus.MustRegister(handler.IdApiSummary)
	prometheus.MustRegister(consumptionCounter)
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := store.InitiateStorageService(); err != nil {
		log.Fatalf("Consumer. Error in initiating storage service. Error: [%v]", err)
	}
	storageSvc = store.GetService()
	defer func(Db *badger.DB) {
		err := Db.Close()
		if err != nil {
			log.Printf("Consumer. Error in closing DB connection. Error: [%v]", err)
		}
	}(storageSvc.Db)

	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.BalanceStrategyRoundRobin}
	if oldest {
		config.Consumer.Offsets.Initial = sarama.OffsetOldest
	}

	// Set up a new Sarama consumer group
	consumer := Consumer{
		ready: make(chan bool),
	}

	client, err := sarama.NewConsumerGroup(strings.Split(brokers, ","), group, config)
	if err != nil {
		log.Panicf("Error creating consumer group client: %v", err)
	}

	// Register http routes
	routes.RegisterRoutes()

	// Prometheus metric
	http.Handle("/metrics", promhttp.Handler())

	// Start http server
	fmt.Printf("Starting server at port 8080...\n")
	go func() {
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatal(err)
		}
	}()

	// Register Prometheus custom metrics
	registerPrometheusMetrics()

	log.Println("Starting a new Sarama consumer...")
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := client.Consume(ctx, strings.Split(topic, ","), &consumer); err != nil {
			log.Panicf("Consumer: Error: %v", err)
		}
		// check if context was cancelled, signaling that the consumer should stop
		if ctx.Err() != nil {
			return
		}
		consumer.ready = make(chan bool)
	}()

	<-consumer.ready // wait till the consumer has been set up
	log.Println("Sarama consumer up and running!...")

	select {
	case <-ctx.Done():
		log.Println("terminating: context cancelled")
	}

	wg.Wait()
	if err = client.Close(); err != nil {
		log.Panicf("Consumer: error closing client: %v", err)
	}
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(consumer.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message := <-claim.Messages():
			if err := processMessage(message); err != nil {
				return err
			}
			session.MarkMessage(message, "")

		case <-session.Context().Done():
			return nil
		}
	}
}

func processMessage(message *sarama.ConsumerMessage) error {
	var consumedMessage consumer_structs.Message
	if err := json.Unmarshal(message.Value, &consumedMessage); err != nil {
		log.Printf("Consumer: Error in formatting consumed message. Error: [%v]", err)
		return err
	}
	log.Printf("Message consumed: [%v]", consumedMessage)

	// Save consumed message in badger KV store
	if err := storageSvc.SaveConsumedMessage(consumedMessage); err != nil {
		log.Printf("Consumer: Error processing consumed message. Error: [%v]", err)
		return err
	}

	// Update consumption counter metric
	consumptionCounter.WithLabelValues(consumedMessage.Id).Inc()

	return nil
}
