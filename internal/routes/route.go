package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/jfosburgh/go-rss/internal/database"
)

type apiConfig struct {
	DB *database.Queries
}

func CreateRouter(dbQueries *database.Queries) chi.Router {
	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{}))

	config := apiConfig{
		DB: dbQueries,
	}

	v1 := createV1Router(&config)
	r.Mount("/v1", v1)

	return r
}
