package kernel

import (
	"net/http"

	"github.com/Fcou/web-frame/framework/gin"
)

// 引擎服务
type FcouKernelService struct {
	engine *gin.Engine
}

// 初始化web引擎服务实例
func NewFcouKernelService(params ...interface{}) (interface{}, error) {
	httpEngine := params[0].(*gin.Engine)
	return &FcouKernelService{engine: httpEngine}, nil
}

// 返回web引擎
func (s *FcouKernelService) HttpEngine() http.Handler {
	return s.engine
}
