package app

import (
	"github.com/Fcou/web-frame/framework"
	"github.com/Fcou/web-frame/framework/contract"
)

// FcouAppProvider 提供App的具体实现方法
type FcouAppProvider struct {
	BaseFolder string //BaseFolder 是获取项目的基础路径,再按照目录设计项目结构
}

// Register 注册FcouApp方法
func (h *FcouAppProvider) Register(container framework.Container) framework.NewInstance {
	return NewFcouApp
}

// Boot 启动调用
func (h *FcouAppProvider) Boot(container framework.Container) error {
	return nil
}

// IsDefer 是否延迟初始化
func (h *FcouAppProvider) IsDefer() bool {
	return false
}

// Params 获取初始化参数
func (h *FcouAppProvider) Params(container framework.Container) []interface{} {
	return []interface{}{container, h.BaseFolder}
}

// Name 获取字符串凭证
func (h *FcouAppProvider) Name() string {
	return contract.AppKey
}
