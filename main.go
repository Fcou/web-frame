// Copyright 2021 jianfengye.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
package main

import (
	"github.com/Fcou/web-frame/app/console"
	"github.com/Fcou/web-frame/app/http"
	"github.com/Fcou/web-frame/framework"
	"github.com/Fcou/web-frame/framework/provider/app"
	"github.com/Fcou/web-frame/framework/provider/cache"
	"github.com/Fcou/web-frame/framework/provider/config"
	"github.com/Fcou/web-frame/framework/provider/distributed"
	"github.com/Fcou/web-frame/framework/provider/env"
	"github.com/Fcou/web-frame/framework/provider/id"
	"github.com/Fcou/web-frame/framework/provider/kernel"
	"github.com/Fcou/web-frame/framework/provider/log"
	"github.com/Fcou/web-frame/framework/provider/orm"
	"github.com/Fcou/web-frame/framework/provider/redis"
	"github.com/Fcou/web-frame/framework/provider/ssh"
	"github.com/Fcou/web-frame/framework/provider/trace"
)

func main() {
	// 初始化服务容器
	container := framework.NewFcouContainer()
	// 绑定App服务提供者
	container.Bind(&app.FcouAppProvider{})
	// 后续初始化需要绑定的服务提供者...
	container.Bind(&env.FcouEnvProvider{})
	container.Bind(&distributed.LocalDistributedProvider{})
	container.Bind(&config.FcouConfigProvider{})
	container.Bind(&id.FcouIDProvider{})
	container.Bind(&trace.FcouTraceProvider{})
	container.Bind(&log.FcouLogServiceProvider{})
	container.Bind(&orm.GormProvider{})
	container.Bind(&redis.RedisProvider{})
	container.Bind(&cache.FcouCacheProvider{})
	container.Bind(&ssh.SSHProvider{})

	// 将HTTP引擎初始化,并且作为服务提供者绑定到服务容器中
	if engine, err := http.NewHttpEngine(container); err == nil {
		container.Bind(&kernel.FcouKernelProvider{HttpEngine: engine})
	}

	// 运行root命令
	console.RunCommand(container)
}
