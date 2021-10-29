package http

import (
	"github.com/Fcou/web-frame/app/http/module/demo"
	"github.com/Fcou/web-frame/framework/gin"
)

func Routes(r *gin.Engine) {

	r.Static("/dist/", "./dist/")

	demo.Register(r)
}
