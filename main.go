package main

import (
	"database/sql" // postgres driver
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/erbatax/chirpy_go/internal/database"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits               atomic.Int32
	db                           *database.Queries
	jwtSecret                    string
	accessTokenExpiresInSeconds  int64
	refreshTokenExpiresInSeconds int64
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Printf("Can't connect to database. %v", err)
		os.Exit(1)
		return
	}
	defer db.Close()

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		fmt.Println("JWT_SECRET environment variable is not set")
		os.Exit(1)
		return
	}

	port := "8080"

	serveMux := http.NewServeMux()
	server := http.Server{
		Addr:    ":" + port,
		Handler: serveMux,
	}
	apiCfg := &apiConfig{
		db:                           database.New(db),
		jwtSecret:                    jwtSecret,
		accessTokenExpiresInSeconds:  3600,           // Default to 1 hour
		refreshTokenExpiresInSeconds: 60 * 24 * 3600, // Default to 60 days
	}

	serveMux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))

	serveMux.HandleFunc("GET /api/healthz", healthzHandler)

	serveMux.HandleFunc("GET /admin/metrics", apiCfg.hitsHandler)
	serveMux.HandleFunc("POST /admin/reset", apiCfg.resetHandler)

	serveMux.HandleFunc("POST /api/validate_chirp", validateChirpHandler)

	serveMux.HandleFunc("POST /api/users", apiCfg.createUserHandler)
	serveMux.HandleFunc("POST /api/login", apiCfg.loginHandler)
	serveMux.HandleFunc("POST /api/refresh", apiCfg.refreshHandler)
	serveMux.HandleFunc("POST /api/revoke", apiCfg.revokeHandler)

	serveMux.HandleFunc("POST /api/chirps", apiCfg.createChirpHandler)
	serveMux.HandleFunc("GET /api/chirps", apiCfg.getManyChirpsHandler)
	serveMux.HandleFunc("GET /api/chirps/{id}", apiCfg.getOneChirpsHandler)

	log.Printf("Starting server on :%s\n", port)
	server.ListenAndServe()
}
