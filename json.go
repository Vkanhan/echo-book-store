package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithError(w http.ResponseWriter, statusCode int, msg string) {
	if statusCode >= 500 {
		log.Printf("Responding with 5XX error: %s", msg)
	}

	type errResponse struct {
		Error string `json:"error"`
	}

	respondWithJSON(w, statusCode, errResponse{Error: msg})
}

func respondWithJSON(w http.ResponseWriter, statusCode int, payload any) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marshal JSON response: %v", err)
		w.WriteHeader(statusCode)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(data)
}
