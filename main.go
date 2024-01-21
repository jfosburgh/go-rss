package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jfosburgh/go-rss/internal/database"
	"github.com/jfosburgh/go-rss/internal/routes"
	"github.com/jfosburgh/go-rss/internal/scraper"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()
	port := os.Getenv("PORT")

	dbURL := os.Getenv("DB")
	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Could not connect to database")
	}
	dbQueries := database.New(conn)

	r := routes.CreateRouter(dbQueries)

	server := http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: r,
	}

	go scraper.FetchFeeds(dbQueries, 10, 2)
	log.Fatal(server.ListenAndServe())
}
