package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"webserver/internal/database"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {

	godotenv.Load()
	const port = "8080"

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is missing")
	}

	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("error opening database connection: %s", err)
	}

	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db:             database.New(dbConn),
	}

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
