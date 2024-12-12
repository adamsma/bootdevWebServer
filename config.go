package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
	"webserver/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(resp, req)
	})
}

func (cfg *apiConfig) handlerHits(resp http.ResponseWriter, req *http.Request) {

	resp.Header().Set("Content-Type", "text/html; charset=utf-8")
	resp.WriteHeader(http.StatusOK)

	hitHTML := "<html><body><h1>Welcome, Chirpy Admin</h1>"
	hitHTML += fmt.Sprintf("<p>Chirpy has been visited %d times!</p>", cfg.fileserverHits.Load())
	hitHTML += "</body></html>"
	resp.Write([]byte(hitHTML))

}

func (cfg *apiConfig) handlerReset(resp http.ResponseWriter, req *http.Request) {

	cfg.fileserverHits.Store(0)

	resp.Header().Set("Content-Type", "text/plain; charset=utf-8")
	resp.WriteHeader(http.StatusOK)

	resp.Write([]byte("Hits reset to 0"))

}
