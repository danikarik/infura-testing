package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

// M is synonym for map structure.
type M map[string]interface{}

const (
	// InfuraMainNetWS is a Infura's main net address.
	InfuraMainNetWS = "wss://mainnet.infura.io/ws/v3/9ce23ef47beb48d99c27eda019aed08c"
	// InfuraMainNetTCP is a Infura's main net address.
	InfuraMainNetTCP = "https://mainnet.infura.io/v3/9ce23ef47beb48d99c27eda019aed08c"
)

// NewRPC creates a new RPC client.
func NewRPC(ctx context.Context, url string) (*rpc.Client, error) {
	client, err := rpc.DialWebsocket(ctx, url, "")
	if err != nil {
		return nil, err
	}
	return client, nil
}

func main() {
	ctx := context.Background()

	client, err := ethclient.DialContext(ctx, InfuraMainNetWS)
	if err != nil {
		log.Fatal(err)
	}

	infuraClient, err := NewRPC(ctx, InfuraMainNetWS)
	if err != nil {
		log.Fatal(err)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(SetContentType("application/json"))

	r.Get("/", InfoHandler)

	r.Route("/api/v1.1/{networkID}/transactions", func(r chi.Router) {

		r.Get("/pending", func(w http.ResponseWriter, r *http.Request) {
			ctx := context.Background()
			filterID := r.URL.Query().Get("filterId")
			from := r.URL.Query().Get("from")

			var fromAddress *common.Address
			if from != "" {
				fromAddr := common.HexToAddress(from)
				fromAddress = &fromAddr
			}

			opts := &PendingTransactionsQuery{
				FilterID: filterID,
				From:     fromAddress,
			}

			result, err := PendingTransactions(ctx, infuraClient, client, opts)
			if err != nil {
				json.NewEncoder(w).Encode(M{
					"success": false,
					"message": err.Error(),
				})
				return
			}

			if filterID != "" {
				cacheKey := fmt.Sprintf(ETHPendingSetTemplate, filterID)
				set := make(map[string]struct{})

				data, ok := httpCache.Get(cacheKey)
				if ok {
					set = data.(map[string]struct{})
				}

				n := []*Transaction{}
				for _, tx := range result.Transactions {
					if _, ok := set[tx.Hash().String()]; ok {
						continue
					}
					set[tx.Hash().String()] = struct{}{}
					n = append(n, tx)
				}

				httpCache.Set(cacheKey, set, ETHPedingSetTTL)

				result.Transactions = n
			}

			json.NewEncoder(w).Encode(result)
		})

		r.Get("/filter", func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			var result string
			err := infuraClient.CallContext(ctx, &result, "eth_newPendingTransactionFilter")
			if err != nil {
				json.NewEncoder(w).Encode(M{
					"success": false,
					"message": err.Error(),
				})
				return
			}
			json.NewEncoder(w).Encode(M{
				"success": true,
				"result":  result,
			})
		})

		r.Get("/changes/{filterID}", func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			filterID := chi.URLParam(r, "filterID")
			var result []string
			err := infuraClient.CallContext(ctx, &result, "eth_getFilterChanges", filterID)
			if err != nil {
				json.NewEncoder(w).Encode(M{
					"success": false,
					"message": err.Error(),
				})
				return
			}
			json.NewEncoder(w).Encode(M{
				"success": true,
				"result":  result,
			})
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
