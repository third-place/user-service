package kafka

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"os"
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
	return kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": os.Getenv("KAFKA_BOOTSTRAP_SERVERS"),
		"security.protocol": os.Getenv("KAFKA_SECURITY_PROTOCOL"),
		"sasl.mechanisms":   os.Getenv("KAFKA_SASL_MECHANISM"),
		"sasl.username":     os.Getenv("KAFKA_SASL_USERNAME"),
		"sasl.password":     os.Getenv("KAFKA_SASL_PASSWORD"),
	})
}
