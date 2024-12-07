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

	sMux.HandleFunc("/healthz", handlerHealth)
	sMux.HandleFunc("/metrics", apiCfg.handlerHits)
	sMux.HandleFunc("/reset", apiCfg.handlerReset)

	server := &http.Server{
		Handler: sMux,
		Addr:    ":" + port,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())

}
