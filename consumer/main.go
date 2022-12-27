package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"producer/consumer/consumer_structs"
	"producer/consumer/routes"
	"producer/consumer/store"
	"producer/helper"
	"producer/structs"
	"strings"
	"sync"

	"github.com/Shopify/sarama"
)

// Consumer represents a Sarama consumer group consumer
type Consumer struct {
	ready chan bool
}

// Sarama configuration options
var (
	brokers        = "0.0.0.0:8097"
	topic          = "user_details_1"
	group          = "user_group_1"
	oldest         = true
	commonConfig   structs.CommonConfig
	consumerConfig consumer_structs.ConsumerConfig
	storageSvc     store.StorageService
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// TODO: get filenames from consts
	commonConfig = helper.LoadCommonConfiguration("config/common.json")
	consumerConfig = helper.LoadConsumerConfiguration("consumer/config/config.json")
	storageSvc = store.StorageService{
		ConsumerConfig: consumerConfig,
		CommonConfig:   commonConfig,
	}
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.BalanceStrategyRoundRobin}
	if oldest {
		config.Consumer.Offsets.Initial = sarama.OffsetOldest
	}

	// Setup a new Sarama consumer group
	consumer := Consumer{
		ready: make(chan bool),
	}

	client, err := sarama.NewConsumerGroup(strings.Split(brokers, ","), group, config)
	if err != nil {
		log.Panicf("Error creating consumer group client: %v", err)
	}

	// start http server
	routes.RegisterRoutes()

	fmt.Printf("Starting server at port 8080...\n")
	go func() {
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatal(err)
		}
	}()

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

	<-consumer.ready // Await till the consumer has been set up
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
	var consumedMessage structs.Message
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

	return nil
}
