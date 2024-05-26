package main

import (
	"fmt"
	"github.com/rs/cors"
	"github.com/third-place/user-service/internal"
	"github.com/third-place/user-service/internal/middleware"
	"log"
	"net/http"
	"os"
	"strconv"
)

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

func main() {
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
