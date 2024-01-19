package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func CreateRouter() chi.Router {
	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{}))

	v1 := createV1Router()
	r.Mount("/v1", v1)

	return r
}
