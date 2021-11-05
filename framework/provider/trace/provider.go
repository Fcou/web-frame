package trace

import (
	"github.com/Fcou/web-frame/framework"
	"github.com/Fcou/web-frame/framework/contract"
)

type FcouTraceProvider struct {
	c framework.Container
}

// Register registe a new function for make a service instance
func (provider *FcouTraceProvider) Register(c framework.Container) framework.NewInstance {
	return NewFcouTraceService
}

// Boot will called when the service instantiate
func (provider *FcouTraceProvider) Boot(c framework.Container) error {
	provider.c = c
	return nil
}

// IsDefer define whether the service instantiate when first make or register
func (provider *FcouTraceProvider) IsDefer() bool {
	return false
}

// Params define the necessary params for NewInstance
func (provider *FcouTraceProvider) Params(c framework.Container) []interface{} {
	return []interface{}{provider.c}
}

/// Name define the name for this service
func (provider *FcouTraceProvider) Name() string {
	return contract.TraceKey
}
