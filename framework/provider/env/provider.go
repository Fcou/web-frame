package env

import (
	"github.com/Fcou/web-frame/framework"
	"github.com/Fcou/web-frame/framework/contract"
)

type FcouEnvProvider struct {
	Folder string
}

// Register registe a new function for make a service instance
func (provider *FcouEnvProvider) Register(c framework.Container) framework.NewInstance {
	return NewFcouEnv
}

// Boot will called when the service instantiate
func (provider *FcouEnvProvider) Boot(c framework.Container) error {
	app := c.MustMake(contract.AppKey).(contract.App)
	provider.Folder = app.BaseFolder()
	return nil
}

// IsDefer define whether the service instantiate when first make or register
func (provider *FcouEnvProvider) IsDefer() bool {
	return false
}

// Params define the necessary params for NewInstance
func (provider *FcouEnvProvider) Params(c framework.Container) []interface{} {
	return []interface{}{provider.Folder}
}

/// Name define the name for this service
func (provider *FcouEnvProvider) Name() string {
	return contract.EnvKey
}
