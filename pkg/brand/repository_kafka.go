package brand

import (
	"encoding/json"
	"fmt"

	"github.com/markus-azer/products-service/pkg/entity"

	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

//KafkaRepository mongodb repo
type KafkaRepository struct {
	consumer *kafka.Consumer
}

//NewKafkaRepository create new repository
func NewKafkaRepository(c *kafka.Consumer) MessagesRepository {
	return &KafkaRepository{
		consumer: c,
	}
}

//GetMessages pull messages from kafka brand topic
func (r *KafkaRepository) GetMessages() <-chan entity.Message {
	r.consumer.SubscribeTopics([]string{"brands"}, nil)
	c := make(chan entity.Message)
	go func() {
		for {
			msg, err := r.consumer.ReadMessage(-1)

			if err == nil {
				var m entity.Message
				json.Unmarshal(msg.Value, &m)
				m.ID = string(msg.Key)
				//c <- entity.Message{ID: string(msg.Key), Version: 1, Type: "BRAND_CREATED", Payload: string(msg.Value), Timestamp: msg.Timestamp}
				c <- m
				fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))

			} else {
				// The client will automatically try to recover from all errors.
				fmt.Printf("Consumer error: %v (%v)\n", err, msg)
			}
		}
	}()

	return c
}
