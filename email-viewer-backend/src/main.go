package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type SearchResult struct {
	Results string `json:"results"`
}

type SearchRequest struct {
	Query string `json:"query"`
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(corsMiddleware)

	r.Post("/search", func(res http.ResponseWriter, req *http.Request) {
		var body SearchRequest

		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			http.Error(res, "Invalid request body", http.StatusBadRequest)
			return
		}

		fmt.Println("Received query: " + body.Query)

		zincResponse, err := queryZinc(body.Query, "")

		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}

		res.Header().Set("Content-Type", "application/json")

		err = json.NewEncoder(res).Encode(zincResponse.Hits.Hits)

		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}
	})

	http.ListenAndServe(":3000", r)
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		// Check if it's a preflight request
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
