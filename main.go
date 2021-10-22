package main

import (
	"net/http"
	"web-frame/framework"
	"web-frame/route"
)

func main() {
	core := framework.NewCore()
	route.RegisterRouter(core)
	server := &http.Server{
		Handler: core,
		Addr:    ":8888",
	}
	server.ListenAndServe()
}
