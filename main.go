package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jfosburgh/go-rss/internal/routes"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	port := os.Getenv("PORT")

	r := routes.CreateRouter()

	server := http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: r,
	}

	log.Fatal(server.ListenAndServe())
}
