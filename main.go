package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/samnodier/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
	jwtSecret      string
}

func fileserverHandler(filepathRoot string) http.HandlerFunc {
	fileServer := http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))
	return func(w http.ResponseWriter, r *http.Request) {
		fileServer.ServeHTTP(w, r)
	}
}

func (cfg *apiConfig) countNumReq(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	hits := fmt.Sprintf(`
		<html>
		<body>
			<h1>Welcome, Chirpy Admin</h1>
			<p>Chirpy has been visited %d times!</p>
		</body>
		</html>
	`, cfg.fileserverHits.Load())
	w.Write([]byte(hits))
}

func main() {
	godotenv.Load()
	const filepathRoot = "."
	const port = "8080"
	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")
	jwtSecret := os.Getenv("JWT_SECRET")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Printf("error opening the database: %s", err)
	}
	dbQueries := database.New(db)
	cfg := apiConfig{
		db:        dbQueries,
		platform:  platform,
		jwtSecret: jwtSecret,
	}
	mux := http.NewServeMux()
	mux.Handle("/app/", cfg.middlewareMetricsInc(fileserverHandler((filepathRoot))))
	mux.HandleFunc("GET /api/healthz", handlerReadiness)

	mux.HandleFunc("GET /admin/metrics", cfg.countNumReq)
	mux.HandleFunc("POST /admin/reset", cfg.handlerReset)

	mux.HandleFunc("POST /api/users", cfg.handleUsersCreate)
	mux.HandleFunc("PUT /api/users", cfg.handleUserUpdate)

	mux.HandleFunc("POST /api/login", cfg.handleLogin)

	mux.HandleFunc("GET /api/chirps", cfg.handleChirpsGet)
	mux.HandleFunc("GET /api/chirps/{chirpID}", cfg.handleChirpGet)
	mux.HandleFunc("POST /api/chirps", cfg.handleChirpsCreate)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", cfg.handleChirpDelete)

	mux.HandleFunc("POST /api/refresh", cfg.handleRefresh)

	mux.HandleFunc("POST /api/revoke", cfg.handleRevoke)

	mux.HandleFunc("POST /api/polka/webhooks", cfg.handleUserUpgrade)

	server := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s...\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}
