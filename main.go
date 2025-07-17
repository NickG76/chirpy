package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
	"os"
	"database/sql"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	gotdotenv.Load()
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
	}

	mux := http.NewServeMux()
	mux.Handle("GET /app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("GET /admin/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	mux.HandleFunc("/api/validate_chirp", apiCfg.validateChirpHandler)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	fmt.Println()
	fmt.Println()
	log.Printf("Server running on: http://localhost:%s/app/\n", port)
	log.Fatal(srv.ListenAndServe())
}
