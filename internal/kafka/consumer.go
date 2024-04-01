package kafka

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/google/uuid"
	"github.com/third-place/user-service/internal/db"
	"github.com/third-place/user-service/internal/model"
	"github.com/third-place/user-service/internal/repository"
	"log"
)

func InitializeAndRunLoop() {
	userRepository := repository.CreateUserRepository(db.CreateDefaultConnection())
	err := loopConsumer(userRepository)
	if err != nil {
		log.Fatal(err)
	}
}

func createConsumer() *kafka.Consumer {
	cfg := createConnectionConfig()
	_ = cfg.SetKey("group.id", "user-service")
	_ = cfg.SetKey("auto.offset.reset", "earliest")
	c, err := kafka.NewConsumer(cfg)

	if err != nil {
		log.Fatal(err)
	}

	_ = c.SubscribeTopics([]string{"images"}, nil)
	return c
}

func loopConsumer(userRepository *repository.UserRepository) error {
	consumer := createConsumer()
	for {
		log.Print("kafka ready to consume messages")
		data, err := consumer.ReadMessage(-1)
		if err != nil {
			log.Print(err)
			return nil
		}
		log.Print("consuming message ", string(data.Value))
		image, err := model.DecodeMessageToImage(data.Value)
		if err != nil {
			log.Print("error decoding message to image, skipping", string(data.Value))
			continue
		}
		userEntity, err := userRepository.GetUserFromUuid(uuid.MustParse(image.User.Uuid))
		if err != nil {
			log.Print("user not found when updating profile pic")
			continue
		}
		log.Print("update user with key", userEntity.Uuid.String(), image.Key)
		userEntity.ProfilePic = image.Key
		userRepository.Save(userEntity)
	}
}
