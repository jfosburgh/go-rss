package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jfosburgh/go-rss/internal/database"
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

func (cfg *apiConfig) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	type params struct {
		Name string `json:"name"`
	}

	jsonDecoder := json.NewDecoder(r.Body)
	userParams := params{}
	err := jsonDecoder.Decode(&userParams)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Could not parse body as user: %v", err))
		return
	}

	id := uuid.New()
	now := time.Now().UTC()
	user, err := cfg.DB.CreateUser(r.Context(), database.CreateUserParams{ID: id, CreatedAt: now, UpdatedAt: now, Name: userParams.Name})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Could not create user: %v", err))
		return
	}

	respondWithJSON(w, 201, user)
}

func (cfg *apiConfig) handleGetUserByAPIKey(w http.ResponseWriter, r *http.Request) {
	apiKey, ok := strings.CutPrefix(r.Header.Get("Authorization"), "ApiKey ")
	if !ok {
		respondWithError(w, 401, "ApiKey required")
		return
	}

	user, err := cfg.DB.GetUserByAPIKey(r.Context(), apiKey)
	if err != nil {
		respondWithError(w, 404, fmt.Sprintf("No user found: %v", err))
		return
	}

	respondWithJSON(w, 200, user)
}

func createV1Router(config *apiConfig) chi.Router {
	v1 := chi.NewRouter()

	v1.Get("/readiness", handleGetReadiness)
	v1.Get("/err", handleGetError)
	v1.Get("/users", config.handleGetUserByAPIKey)

	v1.Post("/users", config.handleCreateUser)

	return v1
}
