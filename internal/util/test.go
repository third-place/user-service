package util

import (
	"context"
	"github.com/goombaio/namegenerator"
	"github.com/jinzhu/gorm"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/third-place/user-service/internal/db"
	"github.com/third-place/user-service/internal/entity"
	"math/rand"
	"os"
	"strconv"
	"time"
)

var NameGenerator namegenerator.Generator
var dbConn *gorm.DB

func init() {
	seed := time.Now().UTC().UnixNano()
	NameGenerator = namegenerator.NewNameGenerator(seed)
}

func RandomUsername() string {
	return NameGenerator.Generate()
}

func RandomEmailAddress() string {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	num := r.Intn(100000)
	return os.Getenv("EMAIL_PREFIX") + "+" + strconv.Itoa(num) + "@" + os.Getenv("EMAIL_DOMAIN")
}

func SetupTestDatabase() *gorm.DB {
	if dbConn != nil {
		return dbConn
	}
	// 1. Create PostgreSQL container request
	containerReq := testcontainers.ContainerRequest{
		Image:        "postgres:latest",
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor:   wait.ForListeningPort("5432/tcp"),
		Env: map[string]string{
			"POSTGRES_DB":       "testdb",
			"POSTGRES_PASSWORD": "postgres",
			"POSTGRES_USER":     "postgres",
		},
	}

	// 2. Start PostgreSQL container
	dbContainer, _ := testcontainers.GenericContainer(
		context.Background(),
		testcontainers.GenericContainerRequest{
			ContainerRequest: containerReq,
			Started:          true,
		})

	// 3.1 Get host and port of PostgreSQL container
	host, _ := dbContainer.Host(context.Background())
	port, _ := dbContainer.MappedPort(context.Background(), "5432")

	// 3.2 Create db connection string and connect
	dbConn = db.CreateConnection(
		host,
		port.Port(),
		"testdb",
		"postgres",
		"postgres",
	)

	migrateDb(dbConn)

	return dbConn
}

func migrateDb(dbConn *gorm.DB) {
	dbConn.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")
	dbConn.AutoMigrate(&entity.User{}, &entity.Invite{})
}
