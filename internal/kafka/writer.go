package kafka

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"os"
)

func CreateWriter() *kafka.Producer {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": os.Getenv("KAFKA_BOOTSTRAP_SERVERS"),
		"security.protocol": os.Getenv("KAFKA_SECURITY_PROTOCOL"),
		"sasl.mechanisms": os.Getenv("KAFKA_SASL_MECHANISM"),
		"sasl.username": os.Getenv("KAFKA_SASL_USERNAME"),
		"sasl.password": os.Getenv("KAFKA_SASL_PASSWORD"),
	})
	if err != nil {
		panic(fmt.Sprintf("Failed to create producer: %s", err))
	}
	return producer
}
