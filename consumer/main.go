package main

import (
	"context"
	"encoding/json"
	"log"
	"strings"
	"sync"

	"github.com/Shopify/sarama"
)

// Consumer represents a Sarama consumer group consumer
type Consumer struct {
	ready chan bool
}

type ForwardPayload struct {
	Topic    string      `json:"topic"`
	DateTime string      `json:"datetime"`
	Value    interface{} `json:"value"`
}

type User struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
	Email string `json:"email"`
}

// Sarama configuration options
var (
	brokers = "0.0.0.0:8097"
	topic   = "user_details_1"
	group   = "user_group_1"
	oldest  = true
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Println("Starting a new Sarama consumer...")
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
	var userValue User
	if err := json.Unmarshal(message.Value, &userValue); err != nil {
		log.Printf("Consumer: Error in formatting received message. Error: [%v]", err)
		return err
	}

	payload := ForwardPayload{
		Topic:    message.Topic,
		DateTime: message.Timestamp.String(),
		Value:    userValue,
	}

	if err := postMessage(payload); err != nil {
		log.Printf("Consumer. Error in sending msg in `abc` api. Error: [%v]", err)
		return err
	}

	return nil
}

func postMessage(payload ForwardPayload) error {
	// Send the payload in a POST request to `abc` api.
	log.Printf("Here for post message! Message: [%v]", payload)
	return nil
}
