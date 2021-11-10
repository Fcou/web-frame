package main

import (
	"github.com/Fcou/web-frame/app/console"
	"github.com/Fcou/web-frame/app/http"
	"github.com/Fcou/web-frame/framework"
	"github.com/Fcou/web-frame/framework/provider/app"
	"github.com/Fcou/web-frame/framework/provider/distributed"
	"github.com/Fcou/web-frame/framework/provider/env"
	"github.com/Fcou/web-frame/framework/provider/kernel"
)

func main() {
	// 初始化服务容器
	container := framework.NewFcouContainer()
	// 绑定App服务提供者
	container.Bind(&app.FcouAppProvider{})
	// 绑定抢本地分布式锁服务
	container.Bind(&distributed.LocalDistributedProvider{})
	// 绑定环境变量服务（目前有问题）
	container.Bind(&env.FcouEnvProvider{})
	// 绑定配置服务
	// container.Bind(&config.FcouConfigProvider{})
	// 绑定日志服务
	// container.Bind(&id.FcouIDProvider{})
	// container.Bind(&trace.FcouTraceProvider{})
	// container.Bind(&log.FcouLogServiceProvider{})

	// 将HTTP引擎初始化,并且作为服务提供者绑定到服务容器中
	if engine, err := http.NewHttpEngine(); err == nil {
		container.Bind(&kernel.FcouKernelProvider{HttpEngine: engine})
	}

	// 运行root命令
	console.RunCommand(container)

}
