package routes

import (
	"encoding/json"
	"net/http"
)

func respondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	dat, _ := json.Marshal(payload)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(dat)

	return
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	type errorvals struct {
		Body string `json:"error"`
	}

	respondWithJSON(w, code, errorvals{Body: msg})
}
