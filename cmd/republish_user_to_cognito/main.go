package main

import (
	"github.com/joho/godotenv"
	"github.com/third-place/user-service/internal/db"
	"github.com/third-place/user-service/internal/repository"
	"github.com/third-place/user-service/internal/service"
	"log"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	username := os.Args[1]
	password := os.Args[2]
	userService := service.CreateDefaultUserService()
	userRepository := repository.CreateUserRepository(db.CreateDefaultConnection())
	user, err := userRepository.GetUserFromUsername(username)
	if err != nil {
		log.Print(err)
		return
	}
	log.Printf("signing up with username :: %s, password :: %s", user.Email, password)
	response, err := userService.PublishToCognito(user, password)
	if err != nil {
		log.Print("error :: ", err)
		return
	}
	log.Print(response.String())
}
