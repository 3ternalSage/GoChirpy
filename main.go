package main

import (
	"fmt"
	"net/http"
)

type Handler interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
}

func main() {
	var serverMux http.ServeMux = *http.NewServeMux()
	var server http.Server

	serverMux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(200)
		w.Write([]byte("OK"))
	})
	serverMux.Handle("/app/*", http.StripPrefix("/app", http.FileServer(http.Dir("."))))

	server.Handler = &serverMux
	server.Addr = "localhost:8080"
	ok := server.ListenAndServe()
	fmt.Println(ok)
}
