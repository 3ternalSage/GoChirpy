package main

import (
	"fmt"
	"net/http"
)

func main() {
	var serverMux http.ServeMux = *http.NewServeMux()
	var server http.Server
	server.Handler = &serverMux
	server.Addr = "localhost:8080"
	ok := server.ListenAndServe()
	fmt.Println(ok)
}
