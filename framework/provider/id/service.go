package id

import (
	"github.com/rs/xid"
)

type FcouIDService struct {
}

func NewFcouIDService(params ...interface{}) (interface{}, error) {
	return &FcouIDService{}, nil
}

func (s *FcouIDService) NewID() string {
	return xid.New().String()
}
