package main

import (
	"github.com/Fcou/web-frame/app/console"
	"github.com/Fcou/web-frame/app/http"
	"github.com/Fcou/web-frame/framework"
	"github.com/Fcou/web-frame/framework/provider/app"
	"github.com/Fcou/web-frame/framework/provider/kernel"
)

func main() {
	// 初始化服务容器
	container := framework.NewFcouContainer()
	// 绑定App服务提供者
	container.Bind(&app.FcouAppProvider{})
	// 后续初始化需要绑定的服务提供者...

	// 将HTTP引擎初始化,并且作为服务提供者绑定到服务容器中
	if engine, err := http.NewHttpEngine(); err == nil {
		container.Bind(&kernel.FcouKernelProvider{HttpEngine: engine})
	}

	// 运行root命令
	console.RunCommand(container)

}
