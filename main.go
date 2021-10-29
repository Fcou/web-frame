package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	fcouHttp "github.com/Fcou/web-frame/app/http"
	"github.com/Fcou/web-frame/app/provider/demo"
	"github.com/Fcou/web-frame/framework/middleware"
	"github.com/Fcou/web-frame/framework/provider/app"

	"github.com/Fcou/web-frame/framework/gin"
)

func main() {
	// 创建engine结构
	core := gin.New()
	// 绑定具体的服务
	core.Bind(&app.HadeAppProvider{})
	core.Bind(&demo.DemoProvider{})

	core.Use(gin.Recovery())
	core.Use(middleware.Cost())

	fcouHttp.Routes(core)

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
		log.Fatal("Server Shutdown:", err)
	}
}
