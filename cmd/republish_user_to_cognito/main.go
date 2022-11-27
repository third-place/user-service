package main

import (
	"github.com/danielmunro/otto-user-service/internal/db"
	"github.com/danielmunro/otto-user-service/internal/repository"
	"github.com/danielmunro/otto-user-service/internal/service"
	"github.com/joho/godotenv"
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
	log.Printf("signing up with username :: %s, password :: %s", user.CurrentEmail, password)
	response, err := userService.PublishToCognito(user, password)
	if err != nil {
		log.Print("error :: ", err)
		return
	}
	log.Print(response.String())
}
