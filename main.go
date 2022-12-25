/*
 * User service
 */

package main

import (
	_ "github.com/joho/godotenv/autoload"
	"github.com/rs/cors"
	"github.com/third-place/user-service/internal"
	"github.com/third-place/user-service/internal/kafka"
	"github.com/third-place/user-service/internal/middleware"
	"log"
	"net/http"
)

func main() {
	go readKafka()
	serveHttp()
}

func serveHttp() {
	log.Print("Listening on 8080")
	router := internal.NewRouter()
	handler := cors.AllowAll().Handler(router)
	log.Fatal(
		http.ListenAndServe(
			":8080",
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
