/*
 * Otto user service
 */

package main

import (
	"github.com/danielmunro/otto-user-service/internal"
	"github.com/danielmunro/otto-user-service/internal/middleware"
	_ "github.com/joho/godotenv/autoload"
	"github.com/rs/cors"
	"log"
	"net/http"
)

func main() {
	serveHttp()
}

func serveHttp() {
	log.Print("Listening on 8080")
	router := internal.NewRouter()
	handler := cors.AllowAll().Handler(router)
	log.Fatal(http.ListenAndServe(":8080",
		middleware.CorsMiddleware(middleware.ContentTypeMiddleware(handler))))
}
