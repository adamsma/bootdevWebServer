package main

import (
	"log"
	"net/http"
)

func main() {

	const filepathRoot = "."
	const port = "8080"

	sMux := http.NewServeMux()
	sMux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir("."))))

	sMux.HandleFunc("/healthz", handlerHealth)

	server := &http.Server{
		Handler: sMux,
		Addr:    ":" + port,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())

}

func handlerHealth(resp http.ResponseWriter, req *http.Request) {

	resp.Header().Set("Content-Type", "text/plain; charset=utf-8")
	resp.WriteHeader(http.StatusOK)

	resp.Write([]byte(http.StatusText(http.StatusOK)))

}
