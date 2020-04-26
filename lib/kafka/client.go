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
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "host1:9095,host2:9096,host3:9097"})
	c, err := kafka.NewConsumer(&kafka.ConfigMap{"bootstrap.servers": "host1:9095,host2:9096,host3:9097", "group.id": "products-consumer"})
	if err != nil {
		return nil, err
	}

	fmt.Printf("Created Producer %v\n", p)
	fmt.Printf("Created Consumer %v\n", c)

	client = new(Client)
	client.Producer = p
	client.Consumer = c

	return client, nil
}
