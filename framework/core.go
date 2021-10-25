package framework

import (
	"log"
	"net/http"
	"strings"
)

// type Handler interface {
//	 ServeHTTP(ResponseWriter, *Request)
// }
// 核心Core
// 第一步：实现Handler接口，也就是实现ServeHttp函数
// 第二步：添加路由功能
// 第三步：添加注册中间件存储结构
type Core struct {
	router      map[string]*Tree
	middlewares []ControllerHandler // 从core这边设置的中间件
}

func NewCore() *Core {
	// 初始化路由
	router := map[string]*Tree{}
	router["GET"] = NewTree()
	router["POST"] = NewTree()
	router["PUT"] = NewTree()
	router["DELETE"] = NewTree()
	return &Core{router: router}
}

// 注册中间件
func (c *Core) Use(middlewares ...ControllerHandler) {
	c.middlewares = middlewares
}

// 匹配GET 方法, 增加路由规则
func (c *Core) Get(url string, handlers ...ControllerHandler) {
	if err := c.router["GET"].AddRouter(url, handlers); err != nil {
		log.Fatal("add router error: ", err)
	}
}

// 匹配POST 方法, 增加路由规则
func (c *Core) Post(url string, handlers ...ControllerHandler) {
	if err := c.router["POST"].AddRouter(url, handlers); err != nil {
		log.Fatal("add router error: ", err)
	}
}

// 匹配PUT 方法, 增加路由规则
func (c *Core) Put(url string, handlers ...ControllerHandler) {
	if err := c.router["PUT"].AddRouter(url, handlers); err != nil {
		log.Fatal("add router error: ", err)
	}
}

// 匹配DELETE 方法, 增加路由规则
func (c *Core) Delete(url string, handlers ...ControllerHandler) {
	if err := c.router["DELETE"].AddRouter(url, handlers); err != nil {
		log.Fatal("add router error: ", err)
	}
}

// ==== http method wrap end

func (c *Core) Group(prefix string) IGroup {
	return NewGroup(c, prefix)
}

// 匹配路由，如果没有匹配到，返回nil
func (c *Core) FindRouteByRequest(request *http.Request) []ControllerHandler {
	// uri 和 method 全部转换为大写，保证大小写不敏感
	uri := request.URL.Path
	method := request.Method
	upperMethod := strings.ToUpper(method)

	// 查找第一层map
	if methodHandlers, ok := c.router[upperMethod]; ok {
		return methodHandlers.FindHandler(uri)
	}
	return nil
}

func (c *Core) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	// 封装自定义context
	ctx := NewContext(request, response)

	// 寻找路由
	handlers := c.FindRouteByRequest(request)
	if handlers == nil {
		// 如果没有找到，这里打印日志
		ctx.Json(404, "not found")
		return
	}
	ctx.SetHandlers(handlers)

	// 调用路由函数，如果返回err 代表存在内部错误，返回500状态码
	if err := ctx.Next(); err != nil {
		ctx.Json(500, "inner error")
		return
	}
}
