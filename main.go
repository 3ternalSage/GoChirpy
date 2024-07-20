package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
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

func respondWithError(w http.ResponseWriter, code int, msg string) {
	type resError struct {
		Err string `json:"error"`
	}

	errResp := resError{
		Err: msg,
	}
	resp, err := json.Marshal(errResp)
	if err != nil {
		fmt.Println("ERROR SENDING ERROR")
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(resp)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	resp, err := json.Marshal(payload)
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(resp)
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

	serverMux.HandleFunc("POST /api/validate_chirp", func(w http.ResponseWriter, r *http.Request) {
		type reqParameters struct {
			Body string `json:"body"`
		}

		type resClean struct {
			Body string `json:"cleaned_body"`
		}

		decoder := json.NewDecoder(r.Body)
		params := reqParameters{}
		err := decoder.Decode(&params)
		if err != nil {
			respondWithError(w, 500, "Something went wrong")
			return
		}
		if len(params.Body) > 140 {
			respondWithError(w, 400, "Chirp is too long")
			return
		}

		clean := params.Body
		clean = strings.ReplaceAll(clean, "kerfuffle", "****")
		clean = strings.ReplaceAll(clean, "sharbert", "****")
		clean = strings.ReplaceAll(clean, "fornax", "****")

		clean = strings.ReplaceAll(clean, "Kerfuffle", "****")
		clean = strings.ReplaceAll(clean, "Sharbert", "****")
		clean = strings.ReplaceAll(clean, "Fornax", "****")

		cleanResp := resClean{
			Body: clean,
		}
		respondWithJSON(w, 200, cleanResp)
	})

	serverMux.HandleFunc("GET /admin/metrics", apiCfg.Report)
	serverMux.HandleFunc("/api/reset", apiCfg.Reset)

	server.Handler = &serverMux
	server.Addr = "localhost:8080"
	ok := server.ListenAndServe()
	fmt.Println(ok)
}
