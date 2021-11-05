package services

import (
	"io"

	"github.com/Fcou/web-frame/framework"
	"github.com/Fcou/web-frame/framework/contract"
)

type FcouCustomLog struct {
	FcouLog
}

func NewFcouCustomLog(params ...interface{}) (interface{}, error) {
	c := params[0].(framework.Container)
	level := params[1].(contract.LogLevel)
	ctxFielder := params[2].(contract.CtxFielder)
	formatter := params[3].(contract.Formatter)
	output := params[4].(io.Writer)

	log := &FcouConsoleLog{}

	log.SetLevel(level)
	log.SetCtxFielder(ctxFielder)
	log.SetFormatter(formatter)

	log.SetOutput(output)
	log.c = c
	return log, nil
}
