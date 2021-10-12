package framework

import "net/http"

// type Handler interface {
//	 ServeHTTP(ResponseWriter, *Request)
// }
// 核心Core就是实现Handler接口，也就是实现ServeHttp函数
type Core struct {
}

func NewCore() *Core {
	return &Core{}
}

func (c *Core) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	// TODO
}
