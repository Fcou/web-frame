package http

import (
	"github.com/Fcou/web-frame/app/http/module/demo"
	"github.com/Fcou/web-frame/framework/gin"
	"github.com/Fcou/web-frame/framework/middleware/static"
)

// Routes 绑定业务层路由
func Routes(r *gin.Engine) {

	// 路径先去./dist目录下查找文件是否存在，找到使用文件服务提供服务
	r.Use(static.Serve("/", static.LocalFile("./fcou/dist", false)))

	// 动态路由定义
	demo.Register(r)
}
