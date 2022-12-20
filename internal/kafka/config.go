package kafka

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"os"
)

func createConnectionConfig() *kafka.ConfigMap {
	cfg := &kafka.ConfigMap{
		"bootstrap.servers": os.Getenv("KAFKA_BOOTSTRAP_SERVERS"),
	}
	if protocol := os.Getenv("KAFKA_SECURITY_PROTOCOL"); protocol != "" {
		_ = cfg.SetKey("security.protocol", protocol)
	}
	if mechanism := os.Getenv("KAFKA_SASL_MECHANISM"); mechanism != "" {
		_ = cfg.SetKey("sasl.mechanisms", mechanism)
	}
	if username := os.Getenv("KAFKA_SASL_USERNAME"); username != "" {
		_ = cfg.SetKey("sasl.username", username)
	}
	if password := os.Getenv("KAFKA_SASL_PASSWORD"); password != "" {
		_ = cfg.SetKey("sasl.password", password)
	}
	return cfg
}
