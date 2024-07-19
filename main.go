package main

import (
	"fmt"
	"net/http"
)

type apiConfig struct {
	fileserverHits int
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits++
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) Report(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(fmt.Sprintf("Hits: %d", cfg.fileserverHits)))
}

func (cfg *apiConfig) Reset(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits = 0
}

func main() {
	var serverMux http.ServeMux = *http.NewServeMux()
	var server http.Server
	var apiCfg apiConfig
	apiCfg.fileserverHits = 0

	serverMux.Handle("/app/*", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))

	serverMux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(200)
		w.Write([]byte("OK"))
	})

	serverMux.HandleFunc("GET /metrics", apiCfg.Report)
	serverMux.HandleFunc("/reset", apiCfg.Reset)

	server.Handler = &serverMux
	server.Addr = "localhost:8080"
	ok := server.ListenAndServe()
	fmt.Println(ok)
}
