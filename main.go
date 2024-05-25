/*
 * User service
 */

package main

import (
	"fmt"
	_ "github.com/joho/godotenv/autoload"
	"github.com/rs/cors"
	"github.com/third-place/user-service/internal"
	"github.com/third-place/user-service/internal/kafka"
	"github.com/third-place/user-service/internal/middleware"
	"log"
	"net/http"
	"os"
	"strconv"
)

func main() {
	go readKafka()
	serveHttp()
}

func getServicePort() int {
	port, ok := os.LookupEnv("SERVICE_PORT")
	if !ok {
		port = "8080"
	}
	servicePort, err := strconv.Atoi(port)
	if err != nil {
		log.Fatal(err)
	}
	return servicePort
}

func serveHttp() {
	router := internal.NewRouter()
	handler := cors.AllowAll().Handler(router)
	port := getServicePort()
	log.Printf("Listening on %d", port)
	log.Fatal(
		http.ListenAndServe(
			fmt.Sprintf(":%d", port),
			middleware.CorsMiddleware(
				middleware.ContentTypeMiddleware(handler),
			),
		),
	)
}

func readKafka() {
	log.Print("connecting to kafka")
	kafka.InitializeAndRunLoop()
	log.Print("exit kafka loop")
}
