package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

// middlewareMetricsInc increments the counter for every request to the fileserver
func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

// handlerMetrics returns the current hit count as html
func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	html := fmt.Sprintf(
		`<html>
  		<body>
    					<h1>Welcome, Chirpy Admin</h1>
    					<p>Chirpy has been visited %d times!</p>
  					</body>
					</html>`, cfg.fileserverHits.Load())
	w.Write([]byte(html))
}

// handlerReset resets the hit counter to 0
func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
}

func main() {
	// Initialize the config struct
	apiCfg := &apiConfig{}

	mux := http.NewServeMux()

	// 1. Healthz endpoint
	mux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// 2. Metrics and Reset endpoints (methods on apiCfg)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)

	// 3. FileServer with Middleware and Prefix Stripping
	// Wrap the file server with the middlewareMetricsInc method
	fileServerHandler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(fileServerHandler))

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// 5. Use the ListenAndServe method to start the server
	log.Printf("Serving on port 8080")
	log.Printf("Health check: http://localhost:8080/healthz")
	log.Printf("App (index.html): http://localhost:8080/app/")

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
