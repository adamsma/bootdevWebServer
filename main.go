package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/adamsma/webserver/internal/database"

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

	platform := os.Getenv("PLATFORM")
	if platform == "" {
		log.Fatal("PLATFORM must be set")
	}

	jwtSignSecret := os.Getenv("JWT_SIGN_SECRET")
	if jwtSignSecret == "" {
		log.Fatal("JWT_SIGN_SECRET must be set")
	}

	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db:             database.New(dbConn),
		env:            platform,
		secret:         jwtSignSecret,
	}

	sMux := http.NewServeMux()
	sMux.Handle(
		"/app/",
		apiCfg.middlewareMetricsInc(
			http.StripPrefix("/app", http.FileServer(http.Dir("."))),
		),
	)

	sMux.HandleFunc("GET /api/healthz", handlerHealth)

	sMux.HandleFunc("POST /api/chirps", apiCfg.handleNewChirp)
	sMux.HandleFunc("GET /api/chirps", apiCfg.handleGetChirps)
	sMux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handleGetChirpByID)

	sMux.HandleFunc("POST /api/users", apiCfg.handleCreateUser)

	sMux.HandleFunc("POST /api/login", apiCfg.handleLogin)

	sMux.HandleFunc("POST /api/refresh", apiCfg.handlerRefreshToken)
	sMux.HandleFunc("POST /api/revoke", apiCfg.handlerRevokeRefresh)

	sMux.HandleFunc("GET /admin/metrics", apiCfg.handlerHits)
	sMux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)

	server := &http.Server{
		Handler: sMux,
		Addr:    ":" + port,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(server.ListenAndServe())

}
