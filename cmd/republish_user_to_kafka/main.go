package main

import (
	"encoding/json"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	"github.com/third-place/user-service/internal/db"
	kafka2 "github.com/third-place/user-service/internal/kafka"
	"github.com/third-place/user-service/internal/repository"
	"log"
	"os"
	"time"
)

func main() {
	_ = godotenv.Load()
	username := os.Args[1]
	kafkaWriter, err := kafka2.CreateProducer()
	if err != nil {
		log.Fatal(err)
	}
	userRepository := repository.CreateUserRepository(db.CreateDefaultConnection())
	user, err := userRepository.GetUserFromUsername(username)
	if err != nil {
		log.Fatal("no user found")
	}
	userData, _ := json.Marshal(user)
	topic := "users"
	println("sending to kafka")
	err = kafkaWriter.Produce(
		&kafka.Message{
			Value: userData,
			TopicPartition: kafka.TopicPartition{Topic: &topic,
				Partition: kafka.PartitionAny},
		},
		nil)
	time.Sleep(2 * time.Second)
	if err != nil {
		println("error writing to kafka :: " + err.Error())
		return
	}
	println("done")
}
