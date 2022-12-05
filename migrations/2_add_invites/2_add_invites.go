package main

import (
	"github.com/joho/godotenv"
	"github.com/third-place/user-service/internal/db"
	"github.com/third-place/user-service/internal/entity"
)

func main() {
	_ = godotenv.Load()
	conn := db.CreateDefaultConnection()
	conn.DropTableIfExists(&entity.Invite{})
	conn.CreateTable(&entity.Invite{})
}
