package kafkaStore

import (
	"fmt"

	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

//Client contains Producer Session
type Client struct {
	Producer *kafka.Producer
	Consumer *kafka.Consumer
}

//NewKafkaClient Function that creates Kafka Client
func NewKafkaClient() (*Client, error) {

	var client *Client
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost:9095,localhost:9096,localhost:9097"})
	c, err := kafka.NewConsumer(&kafka.ConfigMap{"bootstrap.servers": "localhost:9095,localhost:9096,localhost:9097", "group.id": "products-consumer"})
	if err != nil {
		return nil, err
	}

	fmt.Printf("Created Producer %v\n", p)
	fmt.Printf("Created Consumer %v\n", c)

	// Delivery report handler for produced messages
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	client = new(Client)
	client.Producer = p
	client.Consumer = c

	return client, nil
}
