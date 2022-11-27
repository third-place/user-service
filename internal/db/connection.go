package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
	"os"
	"time"
)

var dbConn *gorm.DB

func CreateDefaultConnection() *gorm.DB {
	return CreateConnection(
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DBNAME"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"))
}

func CreateConnection(host string, port string, dbname string, user string, password string) *gorm.DB {
	if dbConn == nil {
		db, err := gorm.Open(
			"postgres",
			fmt.Sprintf(
				"host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
				host,
				port,
				dbname,
				user,
				password))
		if err != nil {
			log.Fatal(err)
		}
		dbConn = db
		sqlConnection := dbConn.DB()
		sqlConnection.SetMaxOpenConns(20)
		sqlConnection.SetMaxIdleConns(5)
		sqlConnection.SetConnMaxLifetime(time.Hour)
	}
	return dbConn
}
