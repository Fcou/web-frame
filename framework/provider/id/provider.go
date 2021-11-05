package id

import (
	"github.com/Fcou/web-frame/framework"
	"github.com/Fcou/web-frame/framework/contract"
)

type FcouIDProvider struct {
}

// Register registe a new function for make a service instance
func (provider *FcouIDProvider) Register(c framework.Container) framework.NewInstance {
	return NewFcouIDService
}

// Boot will called when the service instantiate
func (provider *FcouIDProvider) Boot(c framework.Container) error {
	return nil
}

// IsDefer define whether the service instantiate when first make or register
func (provider *FcouIDProvider) IsDefer() bool {
	return false
}

// Params define the necessary params for NewInstance
func (provider *FcouIDProvider) Params(c framework.Container) []interface{} {
	return []interface{}{}
}

/// Name define the name for this service
func (provider *FcouIDProvider) Name() string {
	return contract.IDKey
}
