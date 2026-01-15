package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type SearchRequest struct {
	Intent  string                 `json:"intent"`
	Filters map[string]interface{} `json:"filters"`
}

type Item struct {
	Name  string `json:"name"`
	Price int    `json:"price"`
	Type  string `json:"type"`
}

func main() {
	_ = godotenv.Load() // carga un .env si existe

	port := getEnv("PORT", "8080")
	env := getEnv("ENV", "dev")

	mux := http.NewServeMux()

	// Enpint test
	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{
			"status":  "success",
			"message": "Hello World",
			"time":    time.Now().Format(time.RFC3339),
		})
	})

	// Logger
	handler := withLogging(mux)

	addr := ":" + port
	log.Printf("Starting server on %s (%s)", addr, env)

	srv := &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}
}

func getEnv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func withLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}
