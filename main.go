package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

// M is synonym for map structure.
type M map[string]interface{}

const (
	// InfuraMainNet is a Infura's main net address.
	InfuraMainNet = "wss://mainnet.infura.io/ws/v3/9ce23ef47beb48d99c27eda019aed08c"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(SetContentType("application/json"))

	r.Get("/", InfoHandler)

	r.Route("/api/v1.1/{networdID}/transactions", func(r chi.Router) {

		// TODO eth_newPendingTransactionFilter, eth_getBlockByNumber (true | false), eth_getFilterChanges

		r.Get("/pending", func(w http.ResponseWriter, r *http.Request) {
			json.NewEncoder(w).Encode(M{
				"message": "not implemented",
			})
		})

		r.Get("/filter", func(w http.ResponseWriter, r *http.Request) {
			// TODO
		})
	})

	srv := &http.Server{
		Addr:         ":" + os.Getenv("PORT"),
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	log.Printf("listening on: %s\n", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}

// SetContentType adds `content-type` header to response.
func SetContentType(contentType string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", contentType)
			next.ServeHTTP(w, r)
		})
	}
}

// InfoHandler returns basic info.
func InfoHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(M{
		"description": "infura-testing",
		"author":      "danikarik",
		"version":     "v0.0.1",
	})
}
