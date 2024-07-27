package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/3ternalSage/GoChirpy/database"
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

type resp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

func validateChirp(w http.ResponseWriter, r *http.Request) (string, error) {
	type reqParameters struct {
		Body string `json:"body"`
	}
	decoder := json.NewDecoder(r.Body)
	params := reqParameters{}
	err := decoder.Decode(&params)
	if err != nil {
		print("Decode error")
		print(err)
		respondWithError(w, 500, "Something went wrong")
		return "", err
	}
	if len(params.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return "", nil
	}

	clean := params.Body
	clean = strings.ReplaceAll(clean, "kerfuffle", "****")
	clean = strings.ReplaceAll(clean, "sharbert", "****")
	clean = strings.ReplaceAll(clean, "fornax", "****")

	clean = strings.ReplaceAll(clean, "Kerfuffle", "****")
	clean = strings.ReplaceAll(clean, "Sharbert", "****")
	clean = strings.ReplaceAll(clean, "Fornax", "****")
	return clean, nil
}

func (cfg *apiConfig) Reset(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits = 0
}

var id int = 1

func main() {
	var databasePath string = "database.json"

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

	serverMux.HandleFunc("POST /api/chirps", func(w http.ResponseWriter, r *http.Request) {

		clean, err := validateChirp(w, r)
		if err != nil {
			return
		}
		response := resp{
			ID:   id,
			Body: clean,
		}
		id++
		db, err := database.NewDB(databasePath)
		if err != nil {
			fmt.Print(err)
			respondWithError(w, 500, "Something went wrong")
			return
		}
		_, err = db.CreateChirp(response.Body, response.ID)
		if err != nil {
			fmt.Print("Creating chirp failed")
			respondWithError(w, 500, "Something went wrong")
			return
		}
		respondWithJSON(w, 201, response)
	})

	serverMux.HandleFunc("GET /api/chirps", func(w http.ResponseWriter, r *http.Request) {
		db, err := database.NewDB(databasePath)
		if err != nil {
			fmt.Print(err)
			respondWithError(w, 500, "Something went wrong")
			return
		}
		chirps, err := db.GetChirps()
		if err != nil {
			fmt.Print(err)
			respondWithError(w, 500, "Something went wrong")
			return
		}
		response, err := json.Marshal(chirps)
		if err != nil {
			fmt.Print(err)
			respondWithError(w, 500, "Something went wrong")
			return
		}
		respondWithJSON(w, 200, response)
	})

	serverMux.HandleFunc("GET /admin/metrics", apiCfg.Report)
	serverMux.HandleFunc("/api/reset", apiCfg.Reset)

	server.Handler = &serverMux
	server.Addr = "localhost:8080"
	ok := server.ListenAndServe()
	fmt.Println(ok)
}
