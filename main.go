/*
 * Otto user service
 */

package main

import (
	_ "github.com/joho/godotenv/autoload"
	"github.com/rs/cors"
	"github.com/third-place/user-service/internal"
	"github.com/third-place/user-service/internal/middleware"
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
