package main

import (
	"encoding/json"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/danielmunro/otto-user-service/internal/db"
	kafka2 "github.com/danielmunro/otto-user-service/internal/kafka"
	"github.com/danielmunro/otto-user-service/internal/mapper"
	"github.com/danielmunro/otto-user-service/internal/repository"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	"log"
	"os"
	"time"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	username := os.Args[1]
	kafkaWriter := kafka2.CreateWriter()
	userRepository := repository.CreateUserRepository(db.CreateDefaultConnection())
	user, err := userRepository.GetUserFromUsername(username)
	if err != nil {
		log.Fatal("no user found")
	}
	println("updating user found")
	user.Role = os.Args[2]
	userRepository.Save(user)
	userModel := mapper.MapUserEntityToModel(user)
	userData, _ := json.Marshal(userModel)
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
