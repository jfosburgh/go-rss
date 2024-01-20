package routes

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/jfosburgh/go-rss/internal/database"
)

type authHandler func(http.ResponseWriter, *http.Request, database.User)

func (cfg *apiConfig) authMiddleware(handler authHandler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

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
		handler(w, r, user)
	})
}
