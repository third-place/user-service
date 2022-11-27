package kafka

import (
	"github.com/third-place/user-service/internal/db"
	"github.com/third-place/user-service/internal/model"
	"github.com/third-place/user-service/internal/repository"
	"github.com/google/uuid"
	"log"
)

func InitializeAndRunLoop() {
	userRepository := repository.CreateUserRepository(db.CreateDefaultConnection())
	err := loopKafkaReader(userRepository)
	if err != nil {
		log.Fatal(err)
	}
}

func loopKafkaReader(userRepository *repository.UserRepository) error {
	reader := GetReader()
	for {
		log.Print("kafka ready to consume messages")
		data, err := reader.ReadMessage(-1)
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
		log.Print("update user with s3 key", userEntity.Uuid.String(), image.S3Key)
		userEntity.ProfilePic = image.S3Key
		userRepository.Save(userEntity)
	}
}
