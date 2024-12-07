package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(resp, req)
	})
}

func (cfg *apiConfig) handlerHits(resp http.ResponseWriter, req *http.Request) {

	resp.Header().Set("Content-Type", "text/plain; charset=utf-8")
	resp.WriteHeader(http.StatusOK)

	hitStr := fmt.Sprintf("Hits: %d", cfg.fileserverHits.Load())
	resp.Write([]byte(hitStr))

}

func (cfg *apiConfig) handlerReset(resp http.ResponseWriter, req *http.Request) {

	cfg.fileserverHits.Store(0)

	resp.Header().Set("Content-Type", "text/plain; charset=utf-8")
	resp.WriteHeader(http.StatusOK)

	resp.Write([]byte("Hits reset to 0"))

}
