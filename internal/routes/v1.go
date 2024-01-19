package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func handleGetReadiness(w http.ResponseWriter, r *http.Request) {
	type okresponse struct {
		Status string `json:"status"`
	}

	respondWithJSON(w, 200, okresponse{Status: "ok"})
}

func handleGetError(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, 500, "Internal Server Error")
}

func createV1Router() chi.Router {
	v1 := chi.NewRouter()

	v1.Get("/readiness", handleGetReadiness)
	v1.Get("/err", handleGetError)

	return v1
}
