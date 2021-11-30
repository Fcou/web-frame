package test

import (
	"github.com/Fcou/web-frame/framework"
	"github.com/Fcou/web-frame/framework/provider/app"
	"github.com/Fcou/web-frame/framework/provider/env"
)

const (
	BasePath = "/Users/yejianfeng/Documents/UGit/coredemo/"
)

func InitBaseContainer() framework.Container {
	// 初始化服务容器
	container := framework.NewHadeContainer()
	// 绑定App服务提供者
	container.Bind(&app.HadeAppProvider{BaseFolder: BasePath})
	// 后续初始化需要绑定的服务提供者...
	container.Bind(&env.HadeTestingEnvProvider{})
	return container
}
