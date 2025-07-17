package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
	"os"
	"database/sql"
	"github.com/joho/godotenv"
	"github.com/nickg76/chirpy/internal/database"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	DB *database.Queries
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env files, do they exist?")
	}
	const filepathRoot = "./app/"
	const port = "8080"

	dburl := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dburl)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()
	dbQueries := database.New(db)

	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		DB: dbQueries,
	}

	mux := http.NewServeMux()
	mux.Handle("GET /app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("GET /admin/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	mux.HandleFunc("/api/validate_chirp", apiCfg.validateChirpHandler)
	mux.HandleFunc("POST /api/users", apiCfg.handleCreateUsr)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	fmt.Println()
	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	fmt.Println()
	log.Printf("Server running on: http://localhost:%s/app/\n", port)
	log.Fatal(srv.ListenAndServe())
}
