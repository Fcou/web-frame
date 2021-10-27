# web-frame
#### 跟着轩脉刃从 0 开始构建 Web 框架
#### “先系统设计，再定义接口，最后具体实现”
---
### 01 分析 net/http，创建 Server 数据结构，自实现Handler接口
- Web Server 的本质
```
实际上就是接收、解析 HTTP 请求传输的文本字符，理解这些文本字符的指令，然后进行计算，再将返回值组织成 HTTP 响应的文本字符，通过 TCP 网络传输回去。
```
- 使用go语言的net/http标准库，阅读学习源码的方法：**库函数 > 结构定义 > 结构函数。**
- net/http库 Server 整个逻辑线大致是:**创建服务 -> 监听请求 -> 创建连接 -> 处理请求**
```
第一层，标准库创建 HTTP 服务是通过创建一个 Server 数据结构完成的；
第二层，Server 数据结构在 for 循环中不断监听每一个连接；
第三层，每个连接默认开启一个 Goroutine 为其服务；
第四、五层，serverHandler 结构代表请求对应的处理逻辑，并且通过这个结构进行具体业务逻辑处理；
第六层，Server 数据结构如果没有设置处理函数 Handler，默认使用 DefaultServerMux 处理请求；
第七层，DefaultServerMux 是使用 map 结构来存储和查找路由规则。
```
- 创建 Server 数据结构，并且在数据结构中创建了自定义的 Handler（Core 数据结构）和监听地址，实现了一个 HTTP 服务。
---
### 02 添加上下文 Context 为请求设置超时时间
- 为了防止雪崩，context 标准库的解决思路是：
**在整个树形逻辑链条中，用上下文控制器 Context，实现每个节点的信息传递和共享。**
```
具体操作是：用 Context 定时器为整个链条设置超时时间，时间一到，结束事件被触发，链条中正在处理的服务逻辑会监听到，从而结束整个逻辑链条，让后续操作不再进行。
```
- 在树形逻辑链条上，一个节点其实有两个角色：
   - 下游树的管理者；
   - 上游树的被管理者，那么就对应需要有两个能力：
```
一个是能让整个下游树结束的能力，也就是函数句柄 CancelFunc；
另外一个是在上游树结束的时候被通知的能力，也就是 Done() 方法。同时因为通知是需要不断监听的，所以 Done() 方法需要通过 channel 作为返回值让使用方进行监听。
例子：主线程创建了一个 1 毫秒结束的定时器 Context，在定时器结束的时候，主线程会通过 Done() 函数收到事件结束通知，然后 **主动** 调用函数句柄 cancelFunc 来通知所有子 Context 结束
```
- 阅读context的源码，利用map[canceler]struct{}结构建立起树

```
func WithCancel(parent Context) (ctx Context, cancel CancelFunc) {
	c := newCancelCtx(parent)  //生成一个子上下文
	propagateCancel(parent, &c)  //当父进程被取消时，propagateCancel安排子上下文被取消。
	return &c, func() { c.cancel(true, Canceled) }
}

func newCancelCtx(parent Context) cancelCtx {
	return cancelCtx{Context: parent}
}

type cancelCtx struct {
	Context

	mu       sync.Mutex            // protects following fields
	done     chan struct{}         // created lazily, closed by first cancel call
	children map[canceler]struct{} // set to nil by the first cancel call
	err      error                 // set to non-nil by the first cancel call
}
```
---
### 03 实现路由功能，建立url与处理函数的关系（建立与使用）
- 抽象理解路由功能，就是建立url与处理函数的对应关系，直接想到的是利用map。
	- Method   Request-URI   HandlerFunction 三个变量的对应关系需要两个map嵌套实现
	- map[string]map[string]func
- 为实现**动态路由匹配**，map只能建立1:1的关系，我们需要使用树这种结构来建立n:m的关系
	- 因为有通配符，在匹配 Request-URI 的时候，请求 URI 的某个字符或者某些字符是动态变化的，无法使用 URI 做为 key 来匹配。
	- 这个问题本质是一个字符串匹配，而字符串匹配，比较通用的高效方法就是字典树，也叫 trie 树
- 链式结构：通过函数返回结构本身实现，搞清楚信息输入、输出，按照需求定义即可
```
type IGroup interface {
	// 实现HttpMethod方法
	Get(string, ControllerHandler)
	Post(string, ControllerHandler)
	Put(string, ControllerHandler)
	Delete(string, ControllerHandler)

	// 实现嵌套group
	Group(string) IGroup
}
```
- 实现路由功能步骤
	- 定义树和节点的数据结构
	- 编写函数：“增加路由规则”
	- 编写函数：“查找路由”
	- 将“增加路由规则”和“查找路由”添加到框架中
---
 ### 04 中间件：提高框架的可拓展性

- 设计一个机制，将非业务逻辑代码抽象出来，封装好，提供接口给控制器使用，这个机制的实现，就是中间件。中间件要实现装饰器效果，也就是要把其他业务函数包裹在其中，实现洋葱效果，而不是简单的顺序执行全部函数。
- 目前框架核心逻辑
···
以 Core 为中心，在 Core 中设置路由 router，实现了 Tree 结构，在 Tree 结构中包含路由节点 node；在注册路由的时候，将对应的业务核心处理逻辑 handler ，放在 node 结构的 handler 属性中。
Core 中的 ServeHttp 方法会创建 Context 数据结构，然后 ServeHttp 方法再根据 Request-URI 查找指定 node，并且将 Context 结构和 node 中的控制器 ControllerHandler 结合起来执行具体的业务逻辑。
···
- 从洋葱模型到流水线模型
```
洋葱模型：
func TimeoutHandler(fun ControllerHandler, d time.Duration) ControllerHandler {
  // 使用函数回调
  return func(c *Context) error {
   //...
    }
}
// 超时控制器参数中ControllerHandler结构已经去掉
func Timeout(d time.Duration) framework.ControllerHandler {
  // 使用函数回调
  return func(c *gin.Context) error {
      //...
    }
}
我们可以将每个中间件构造出来的 ControllerHandler 和最终的业务逻辑的 ControllerHandler 结合在一起，都是同样的结构，成为一个 ControllerHandler 数组，也就是控制器链。在最终执行业务代码的时候，能一个个调用控制器链路上的控制器。
```
- 框架，如果提供程序（服务）的注册、使用两个部分，这样会很让用户感到很灵活。
	- node节点上，存储控制器+中间件 数组，是**存储**目的。
	- context上，存储控制器+中间件 数组，是**调用**目的。
	- core上，存储控制器+中间件 数组，是**注册**目的。
---
### 05 封装：让框架更好用
* 定义接口让封装更明确：
	* 请求
		* 参数信息
		* header信息
	* 返回
		* header设置
		* body设置
* **cast库** 实现了多种常见类型之间的相互转换 "github.com/spf13/cast"
* 获取请求地址 url 中带的参数，使用request.URL.Query()
* 获取Form 表单中的参数，使用request.ParseForm()
* 获取请求地址中通配符位置的string，遍历每个node节点，发现是通配符，记录此时的url中的string
* 注意：request.Body 的读取是一次性的，读取一次之后，下个逻辑再去 request.Body 中是读取不到数据内容的。所以我们读取完 request.Body 之后，还要再复制一份 Body 内容，填充到 request.Body 里
```
// 读取文本    
body, err := ioutil.ReadAll(ctx.request.Body)   
if err != nil { 
	return err    
	}    
// 重新填充 request.Body，为后续的逻辑二次读取做准备    
ctx.request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
```
* Json和XML最大的不同是：XML 需要使用 XML 解析器来解析，JSON 可以使用标准的 JavaScript 函数来解析。
* JSONP 是一种我们常用的解决跨域资源共享的方法，获取请求中的参数作为函数名，获取要返回的数据 JSON 作为函数参数，将函数名 + 函数参数作为返回文本
* HTML 输出方法实现，输出 HTML 页面内容的时候，常用“模版 + 数据”的方式。
	* 先根据模版创造出 template 结构；再使用 template.Execute 将传入数据和模版结合。
```
模板文件：
<h1>{{.PageTitle}}</h1>
<ul>
    {{range .Todos}}
        {{if .Done}}
            <li class="done">{{.Title}}</li>
        {{else}}
            <li>{{.Title}}</li>
        {{end}}
    {{end}}
</ul>
传入的数据结构为：
data := TodoPageData{
    PageTitle: "My TODO list",
    Todos: []Todo{
        {Title: "Task 1", Done: false},
        {Title: "Task 2", Done: true},
        {Title: "Task 3", Done: true},
    },
}
```
### 06 重启：进行优雅关闭
* 优雅关闭服务，关闭进程的时候，不能暴力关闭进程，而是要等进程中的所有请求都逻辑处理结束后，才关闭进程。
	* “如何控制关闭进程的操作”
	* “如何等待所有逻辑都处理结束”
* 如何控制关闭进程的操作
	* os/signal 库
	```
	func main() {
	// 这个 Goroutine 是启动服务的 Goroutine
	go func() {
		server.ListenAndServe()
	}()

	// 当前的 Goroutine 等待信号量
	quit := make(chan os.Signal)
	// 监控信号：SIGINT, SIGTERM, SIGQUIT
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	// 这里会阻塞当前 Goroutine 等待信号
	<-quit
	}
	使用 Ctrl 或者 kill 命令，它们发送的信号是进入 main 函数的，即只有 main 函数所在的 Goroutine 会接收到，所以必须在 main 函数所在的 Goroutine 监听信号。
	```
* 如何等待所有逻辑都处理结束
	* 为了实现先阻塞，然后等所有连接处理完再结束退出，Shutdown 使用了两层循环。其中：第一层循环是定时无限循环，每过 ticker 的间隔时间，就进入第二层循环；第二层循环会遍历连接中的所有请求，如果已经处理完操作处于 Idle 状态，就关闭连接，直到所有连接都关闭，才返回。
	```
	ticker := time.NewTicker(shutdownPollInterval) // 设置轮询时间
	defer ticker.Stop()
	for {
			// 真正的操作
		if srv.closeIdleConns() && srv.numListeners() == 0 {
		return lnerr
		}
		select {
		case <-ctx.Done(): // 如果ctx有设置超时，有可能触发超时结束
		return ctx.Err()
		case <-ticker.C:  // 如果没有结束，最长等待时间，进行轮询
		}
	}

	func (s *Server) closeIdleConns() bool {
		s.mu.Lock()
		defer s.mu.Unlock()
		quiescent := true
		for c := range s.activeConn {
			st, unixSec := c.getState()
			// Issue 22682: treat StateNew connections as if
			// they're idle if we haven't read the first request's
			// header in over 5 seconds.
			if st == StateNew && unixSec < time.Now().Unix()-5 {
				st = StateIdle
			}
			if st != StateIdle || unixSec == 0 {
				// Assume unixSec == 0 means it's a very new
				// connection, without state set yet.
				quiescent = false
				continue
			}
			c.rwc.Close()
			delete(s.activeConn, c)
		}
		return quiescent
	}
	```
### 07 集成Gin替换已有核心
* 选框架要根据需要来选，只要最需要的，不要最想要的
	* 小项目：beego
	* 大项目：gin
	```
	如果你开发一个运营管理后台，并发量基本在 100 以下，单机使用，开发的团队规模可能就 1～2 个人，那你最应该考虑功能完备性，明显使用 Beego 的收益会远远大于使用 Gin 和 Echo。因为如果你选择 Gin 和 Echo 之后，还会遇到比如应该选用哪种 ORM、应该选用哪种缓存等一系列问题，而这些在功能组件相当全面的 Beego 中早就被定义好了。
	```