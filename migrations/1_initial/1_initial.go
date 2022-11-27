package main

import (
	"github.com/danielmunro/otto-user-service/internal/db"
	"github.com/danielmunro/otto-user-service/internal/entity"
	"github.com/joho/godotenv"
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