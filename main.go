package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"web-frame/framework"
	"web-frame/framework/middleware"
	"web-frame/route"
)

func main() {
	core := framework.NewCore()

	core.Use(middleware.Recovery())
	core.Use(middleware.Cost())

	route.RegisterRouter(core)
	server := &http.Server{
		Handler: core,
		Addr:    ":8888",
	}

	// 这个goroutine是启动服务的goroutine
	go func() {
		server.ListenAndServe()
	}()

	// 当前的goroutine等待信号量
	quit := make(chan os.Signal)
	// 监控信号：SIGINT, SIGTERM, SIGQUIT
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	// 这里会阻塞当前goroutine等待信号
	<-quit

	// 调用Server.Shutdown graceful结束
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(timeoutCtx); err != nil {
		log.Fatal("Server Grace Shutdown:", err)
	}

	select {
	case <-timeoutCtx.Done():
		log.Println("Timeout of 5 seconds, Server Shutdown")
	}
	log.Println("Server Exiting")
}
