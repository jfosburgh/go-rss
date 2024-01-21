package routes

import (
	// "database/sql"
	"encoding/json"
	"fmt"
	"net/http"
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

func (cfg *apiConfig) handleFeedCreate(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}

	params := parameters{}
	jsonDecoder := json.NewDecoder(r.Body)
	err := jsonDecoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Could not parse body as user: %v", err))
	}

	id := uuid.New()
	now := time.Now().UTC()

	feedParams := database.CreateFeedParams{ID: id, CreatedAt: now, UpdatedAt: now, Name: params.Name, Url: params.Url, UserID: user.ID}
	feed, err := cfg.DB.CreateFeed(r.Context(), feedParams)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Could not create feed: %v", err))
		return
	}

	feedFollowParams := database.CreateFeedFollowParams{ID: uuid.New(), CreatedAt: now, UpdatedAt: now, UserID: user.ID, FeedID: feed.ID}
	feedFollow, err := cfg.DB.CreateFeedFollow(r.Context(), feedFollowParams)
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Feed created successfully but couldn't add feed follow: %v", err))
	}

	type returnparams struct {
		Feed       Feed       `json:"feed"`
		FeedFollow FeedFollow `json:"feed_follow"`
	}

	returnParams := returnparams{Feed: databaseFeedToFeed(feed), FeedFollow: databaseFeedFollowToFeedFollow(feedFollow)}

	respondWithJSON(w, 201, returnParams)
}

func (cfg *apiConfig) handleFeedGetAll(w http.ResponseWriter, r *http.Request) {
	feeds, err := cfg.DB.GetAllFeeds(r.Context())
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Error collecting feeds: %v", err))
	}

	respondWithJSON(w, 200, databaseFeedArrayToFeedArray(feeds))
}

func (cfg *apiConfig) handleFeedFollowCreate(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		FeedId string `json:"feed_id"`
	}

	params := parameters{}
	jsonDecoder := json.NewDecoder(r.Body)
	err := jsonDecoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Could not parse body as user: %v", err))
	}

	feedUUID, err := uuid.Parse(params.FeedId)

	id := uuid.New()
	now := time.Now().UTC()

	feedFollowParams := database.CreateFeedFollowParams{ID: id, CreatedAt: now, UpdatedAt: now, UserID: user.ID, FeedID: feedUUID}
	feedFollow, err := cfg.DB.CreateFeedFollow(r.Context(), feedFollowParams)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Could not create feed follow: %v", err))
		return
	}

	respondWithJSON(w, 201, databaseFeedFollowToFeedFollow(feedFollow))
}

func (cfg *apiConfig) handleFeedFollowDelete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "feedFollowID")

	feedFollowID, err := uuid.Parse(id)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing feedFollowID as UUID: %v", err))
	}

	err = cfg.DB.DeleteFeedFollow(r.Context(), feedFollowID)
	w.WriteHeader(204)
}

func (cfg *apiConfig) handleFeedFollowGet(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollows, err := cfg.DB.GetUserFeedFollows(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, 404, fmt.Sprintf("Error looking up feed follows: %v", err))
	}

	respondWithJSON(w, 200, databaseFeedFollowArrayToFeedFollowArray(feedFollows))
}

func (cfg *apiConfig) handleUserCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}

	params := parameters{}
	jsonDecoder := json.NewDecoder(r.Body)
	err := jsonDecoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Could not parse body as user: %v", err))
	}

	id := uuid.New()
	now := time.Now().UTC()
	user, err := cfg.DB.CreateUser(r.Context(), database.CreateUserParams{ID: id, CreatedAt: now, UpdatedAt: now, Name: params.Name})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Could not create user: %v", err))
		return
	}

	respondWithJSON(w, 201, databaseUserToUser(user))
}

func (cfg *apiConfig) handleUserGet(w http.ResponseWriter, r *http.Request, user database.User) {
	respondWithJSON(w, 200, user)
}

func createV1Router(config *apiConfig) chi.Router {
	v1 := chi.NewRouter()

	v1.Get("/readiness", handleGetReadiness)
	v1.Get("/err", handleGetError)

	v1.Get("/users", config.authMiddleware(config.handleUserGet))
	v1.Post("/users", config.handleUserCreate)

	v1.Get("/feeds", config.handleFeedGetAll)
	v1.Post("/feeds", config.authMiddleware(config.handleFeedCreate))

	v1.Get("/feed_follows", config.authMiddleware(config.handleFeedFollowGet))
	v1.Post("/feed_follows", config.authMiddleware(config.handleFeedFollowCreate))
	v1.Delete("/feed_follows/{feedFollowID}", config.handleFeedFollowDelete)

	return v1
}
