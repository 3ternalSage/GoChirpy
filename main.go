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
	w.Header().Add("Content-Type", "text/html")

	html := "" +
		"<html>\n" +
		"<body>\n" +
		"    <h1>Welcome, Chirpy Admin</h1>\n" +
		fmt.Sprintf(
			"    <p>Chirpy has been visited %d times!</p>\n", cfg.fileserverHits) +
		"</body>\n" +
		"</html>\n"

	w.Write([]byte(html))
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

	serverMux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(200)
		w.Write([]byte("OK"))
	})

	serverMux.HandleFunc("GET /admin/metrics", apiCfg.Report)
	serverMux.HandleFunc("/api/reset", apiCfg.Reset)

	server.Handler = &serverMux
	server.Addr = "localhost:8080"
	ok := server.ListenAndServe()
	fmt.Println(ok)
}
