package main

import (
	"framework"
	"net/http"
)

func main() {
	server := &http.Server{
		Handler: framework.NewCore(),
		Addr:    "localhost:8080",
	}
	server.ListenAndServe()
}
