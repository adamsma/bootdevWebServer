package main

import (
	"log"
	"net/http"
	"sync/atomic"
)

func main() {

	const filepathRoot = "."
	const port = "8080"

	apiCfg := apiConfig{fileserverHits: atomic.Int32{}}

	sMux := http.NewServeMux()
	sMux.Handle(
		"/app/",
		apiCfg.middlewareMetricsInc(
			http.StripPrefix("/app", http.FileServer(http.Dir("."))),
		),
	)

	sMux.HandleFunc("GET /api/healthz", handlerHealth)
	sMux.HandleFunc("GET /admin/metrics", apiCfg.handlerHits)
	sMux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	sMux.HandleFunc("POST /api/validate_chirp", handleValidateChirp)

	server := &http.Server{
		Handler: sMux,
		Addr:    ":" + port,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(server.ListenAndServe())

}
