package services

import (
	"os"

	"github.com/Fcou/web-frame/framework"
	"github.com/Fcou/web-frame/framework/contract"
)

// FcouConsoleLog 代表控制台输出
type FcouConsoleLog struct {
	FcouLog
}

// NewFcouConsoleLog 实例化FcouConsoleLog
func NewFcouConsoleLog(params ...interface{}) (interface{}, error) {
	c := params[0].(framework.Container)
	level := params[1].(contract.LogLevel)
	ctxFielder := params[2].(contract.CtxFielder)
	formatter := params[3].(contract.Formatter)

	log := &FcouConsoleLog{}

	log.SetLevel(level)
	log.SetCtxFielder(ctxFielder)
	log.SetFormatter(formatter)

	// 最重要的将内容输出到控制台
	log.SetOutput(os.Stdout)
	log.c = c
	return log, nil
}
