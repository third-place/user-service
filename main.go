/*
 * Otto user service
 */

package main

import (
	_ "github.com/joho/godotenv/autoload"
	"github.com/rs/cors"
	"github.com/third-place/user-service/internal"
	"github.com/third-place/user-service/internal/kafka"
	"github.com/third-place/user-service/internal/middleware"
	"github.com/third-place/user-service/internal/util"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

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
			util.SessionManager.LoadAndSave(
				middleware.CorsMiddleware(
					middleware.ContentTypeMiddleware(handler),
				),
			),
		),
	)
}

func readKafka() {
	log.Print("connecting to kafka")
	kafka.InitializeAndRunLoop()
	log.Print("exit kafka loop")
}
