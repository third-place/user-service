package kafka

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Producer interface {
	Produce(msg *kafka.Message, deliveryChan chan kafka.Event) error
}

func CreateMessage(data []byte, topic string) *kafka.Message {
	return &kafka.Message{
		Value: data,
		TopicPartition: kafka.TopicPartition{Topic: &topic,
			Partition: kafka.PartitionAny},
	}
}

func CreateProducer() (Producer, error) {
	return kafka.NewProducer(createConnectionConfig())
}
