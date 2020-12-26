package variant

import (
	"bytes"
	"encoding/json"

	"github.com/markus-azer/products-service/pkg/entity"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

//KafkaRepository KafkaRepository
type KafkaRepository struct {
	producer *kafka.Producer
}

//NewKafkaRepository create new repository
func NewKafkaRepository(p *kafka.Producer) MessagesRepository {
	return &KafkaRepository{
		producer: p,
	}
}

//SendMessage Publish new message to kafka
func (r *KafkaRepository) SendMessage(m *entity.Message) {
	// Produce messages to topic (asynchronously)
	topic := "products"
	reqBodyBytes := new(bytes.Buffer)
	json.NewEncoder(reqBodyBytes).Encode(m)

	//TODO: return once msg published
	r.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            []byte(m.ID),
		Value:          []byte(reqBodyBytes.Bytes()),
	}, nil)
}

//SendMessages Publish new messagea to kafka
func (r *KafkaRepository) SendMessages(messages []*entity.Message) {
	// Produce messages to topic (asynchronously)
	topic := "products"

	for _, m := range messages {
		reqBodyBytes := new(bytes.Buffer)
		json.NewEncoder(reqBodyBytes).Encode(m)
		r.producer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Key:            []byte(m.ID),
			Value:          []byte(reqBodyBytes.Bytes()),
		}, nil)
	}
}
