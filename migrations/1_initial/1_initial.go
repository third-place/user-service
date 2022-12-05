package main

import (
	"github.com/joho/godotenv"
	"github.com/third-place/user-service/internal/db"
	"github.com/third-place/user-service/internal/entity"
)

func main() {
	_ = godotenv.Load()
	conn := db.CreateDefaultConnection()
	conn.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\" WITH SCHEMA public;")
	conn.DropTableIfExists(&entity.User{})
	conn.DropTableIfExists(&entity.Password{})
	conn.DropTableIfExists(&entity.Email{})
	conn.CreateTable(&entity.User{})
	conn.CreateTable(&entity.Password{})
	conn.CreateTable(&entity.Email{})
}
