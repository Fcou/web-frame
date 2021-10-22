package main

import (
	"fmt"
	"net/http"
	"web-frame/framework"
)

func main() {
	server := &http.Server{
		Handler: framework.NewCore(),
		Addr:    "localhost:8080",
	}
	server.ListenAndServe()
	fmt.Println("XXXXXXX")
}
