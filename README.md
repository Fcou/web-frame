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
### 08 面向接口编程
* 接口实现了对业务逻辑的抽象，设计接口就是抽象业务的过程。
* 面向接口 / 对象 / 过程
	* “面向过程编程”是指进行业务抽象的时候，我们定义一个一个的过程方法，通过这些过程方法的串联完成具体的业务。
	* “面向对象编程”表示的是在业务抽象的时候，我们先定义业务中的对象，通过这些对象之间的关联来表示整个业务。
	* "面向接口编程"就是对"面向对象编程"的进一步抽象，面对业务，我们并不先定义具体的对象、思考对象有哪些属性，而是先思考如何抽象接口，把接口的定义放在第一步，然后多个模块之间梳理如何通过接口进行交互，最后才是实现具体的模块。
* 按照面向接口编程的理念，将每个模块看成是一个服务。每个模块服务都做两件事情：
	* 一是它和自己提供的接口协议做绑定，这样当其他人要使用这个接口协议时能找到自己。
	* 二是它使用到其他接口协议的时候，去框架主体中寻找。
* 每个模块服务都是一个“服务提供者”（service provider），而我们主体框架需要承担起来的角色叫做“服务容器”（service container），服务容器中绑定了多个接口协议，每个接口协议都由一个服务提供者提供服务。
	* 服务提供者提供的是“创建服务实例的方法”
		* 获取服务凭证的能力 Name；
		* 创建服务实例化方法的能力 Register；
		* 获取服务实例化方法参数的能力 Params；
		* 两个与实例化控制相关的方法，控制实例化时机方法 IsDefer、实例化预处理的方法 Boot。
		```
		// Name 代表了这个服务提供者的凭证
		Name() string
		// NewInstance 定义了如何创建一个新实例，所有服务容器的创建服务
		type NewInstance func(...interface{}) (interface{}, error)
		返回值返回的 interface{} 结构代表了具体的**服务实例**
		// Register 在服务容器中注册了一个实例化服务的方法，是否在注册的时候就实例化这个服务，需要参考 IsDefer 接口。
		Register(Container) NewInstance
		在容器中注册，需要**修改容器**，所以参数传入容器
		服务注册完成后，要能够**创建服务**，NewInstance就是创建服务的方法
		// Params params 定义传递给 NewInstance 的参数，可以自定义多个，建议将 container 作为第一个参数
		Params(Container) []interface{}
		“创建服务实例的方法”的能力，除了实现 NewInstance 方法之外，还需要注册 NewInstance 方法的参数，即可变的 interface{}参数。所以我们的服务提供者还需要提供一个获取服务参数的能力	
		// IsDefer 决定是否在注册的时候实例化这个服务，如果不是注册的时候实例化，那就是在第一次 make 的时候进行实例化操作
		// false 表示不需要延迟实例化，在注册的时候就实例化。true 表示延迟实例化
		IsDefer() bool
		// Boot 在调用实例化服务的时候会调用，可以把一些准备工作：基础配置，初始化参数的操作放在这个里面。
		// 如果 Boot 返回 error，整个服务实例化就会实例化失败，返回错误
		Boot(Container) error
		```
	* 服务容器提供的是“实例化服务的方法”
		* 为服务提供注册绑定
		```
		// Bind 绑定一个服务提供者，如果关键字凭证已经存在，会进行替换操作，不返回 error
		Bind(provider ServiceProvider) error
		```
		* 提供获取服务实例
		```
		// Make 根据关键字凭证获取一个服务
		Make(key string) (interface{}, error)
		```
	* 什么叫**服务实例化**
		* 服务是服务，也就是定义；实例是实例，也就是现在可用的变量。
		* 也就是实例化一个服务struct，创建一个实例，也就是在内存中创建可以使用的变量体
		* instance, err := method(params...)
* 容器和框架的结合
	* 绑定服务操作是全局的操作，将服务容器存放在 Engine 中。
	* 获取服务操作是在**单个请求**中使用的，在 Engine 初始化 Context 的时候，将服务容器传递进入 Context。
	* 接下来完成服务容器方法的封装。
		* Engine 封装 Bind 和 IsBind 方法。(封装也就是调用实际容器的方法)
		* Context 封装 Make、MakeNew、MustMake 方法。
* 如何创建一个服务提供方
	* 要有一个服务接口文件 contract.go，存放服务的接口文件和服务凭证。
	* 需要设计一个 provider.go，这个文件存放服务提供方 ServiceProvider 的实现。
	* 最后在 service.go 文件中实现具体的服务实例。
* 如何通过服务提供方创建服务
	* 需要做两个操作，绑定服务提供方、获取服务
### 09 结构：系统设计框架的整体目录
* 业务代码的目录结构是一种工程化的规范
	* **app**          一个业务就是一个 App
		* console 所有的控制台进程逻辑
		* http 所有Web服务的逻辑, module 目录下每个子目录代表一个模块服务。Web 服务特有的通用中间件，我们使用 http 目录下的 middleware 目录来保存。
		* provider 存放业务提供的服务，每个子目录就代表一个业务服务
	* **config**  存放的是配置文件
	* **framework**
		* command 提供的是框架自带的命令行工具
		* contract 存放框架默认提供的服务协议
		* middleware 存放框架为 Web 服务提供的中间件
		* provider 对应服务协议的具体实现以及服务提供者
		* util 在框架研发过程中通用的一些函数
	* **storage** 存放应用运行过程中产生的内容
		* log 日志
		* runtime 运行的进程 ID 等信息
	* **test**  测试相关的信息
	* **main.go**
	* **README.md**
* 目录结构也是一个服务，其他服务想要使用目录结构的时候，可以通过服务容器，来获取目录结构服务实例
### 10 交互：可以执行命令行的框架才是好框架
* 利用第三方命令行工具库**cobra**
```
// Command代表执行命令的结构
type Command struct {
    ....    
}
// InitFoo 初始化 Foo 命令
func InitFoo() *cobra.Command {
   FooCommand.AddCommand(Foo1Command) //子命令添加方法
   return FooCommand
}
// FooCommand 代表 Foo 命令
var FooCommand = &cobra.Command{
   Use:     "foo", //Use 代表这个命令的调用关键字
   Short:   "foo 的简要说明",
   Long:    "foo 的长说明",
   Aliases: []string{"fo", "f"},
   Example: "foo 命令的例子",
   //RunE 代表当前命令的真正执行函数
   RunE: func(c *cobra.Command, args []string) error {
      container := c.GetContainer()
      log.Println(container)
      return nil
   },
}
// Foo1Command 代表 Foo 命令的子命令 Foo1
var Foo1Command = &cobra.Command{
   Use:     "foo1",
   Short:   "foo1 的简要说明",
   Long:    "foo1 的长说明",
   Aliases: []string{"fo1", "f1"},
   Example: "foo1 命令的例子",
   RunE: func(c *cobra.Command, args []string) error {
      container := c.GetContainer()
      log.Println(container)
      return nil
   },
}
```
* 如何使用命令行 cobra
	* 首先，要把 cobra 库引入到框架中，采用源码引入的方式。
	```
	我们希望把服务容器嵌入到 Command 结构中，让 Command 在调用执行函数 RunE 时，能从参数中获取到服务容器，这样就能从服务容器中使用之前定义的 Make 系列方法获取出具体的服务实例了。
	在根 Command 中设置服务容器
	我们将 Web 服务的启动逻辑封装为一个 Command 命令，将这个 Command 挂载到根 Command 中，然后根据参数获取到这个 Command 节点，执行这个节点中的 RunE 方法，就能启动 Web 服务了。
	```
	* 利用 appStartCommand 启动一个Web服务
* 核心流程
	* 初始化服务容器，将各种服务绑定到容器中
	* 将HTTP引擎初始化，并且作为服务提供者绑定到服务容器中
		* NewHttpEngine 创建了一个绑定了路由的Web引擎 gin.engine
		* 服务容器通过服务提供者FcouKernelProvider绑定该 gin.engine 
	* 创建根Command，为根Command设置服务容器
	* 在根Command上添加各种命令 AddCommand 
		* DemoCommand 命令的 RunE：显示当前路径
			* 获取根Command上的容器 container := c.GetContainer()
			* 从服务容器中获取app的服务实例
			* 调用app服务实例的方法，appService.BaseFolder()，打印出路径
		* appStartCommand 命令的 RunE：启动一个Web服务
			* 获取根Command上的容器
			* 从服务容器中获取kernel的服务实例
			* 从kernel服务实例中获取引擎core
			* 创建一个Server服务
			* goroutine server.ListenAndServe()
			* 优雅关闭
	* 运行根Command
---
### 11 定时任务：让框架支持分布式定时脚本
* 使用利用第三方命令行工具库**cron**库定时执行命令
```
// 创建一个cron实例
c := cron.New()

// 每整点30分钟执行一次
c.AddFunc("30 * * * *", func() { 
  fmt.Println("Every hour on the half hour") 
})
// 上午3-6点，下午8-11点的30分钟执行
c.AddFunc("30 3-6,20-23 * * *", func() {
  fmt.Println(".. in the range 3-6am, 8-11pm") 
})
// 东京时间4:30执行一次
c.AddFunc("CRON_TZ=Asia/Tokyo 30 04 * * *", func() { 
  fmt.Println("Runs at 04:30 Tokyo time every day") 
})
// 从现在开始每小时执行一次
c.AddFunc("@hourly",      func() { 
  fmt.Println("Every hour, starting an hour from now") 
})
// 从现在开始，每一个半小时执行一次
c.AddFunc("@every 1h30m", func() { 
  fmt.Println("Every hour thirty, starting an hour thirty from now") 
})

// 启动cron
c.Start()

...
// 在cron运行过程中增加任务
c.AddFunc("@daily", func() { fmt.Println("Every day") })
..
// 查看运行中的任务
inspect(c.Entries())
..
// 停止cron的运行，优雅停止，所有正在运行中的任务不会停止。
c.Stop() 
```
* AddCronCommand(时间，命令) 是用来创建一个Cron定时任务，封装了cron.AddFunc(时间，匿名函数)
	* 支持秒级别的定时
	* AddCronCommand 函数中核心要做的，就是将**Command 结构的执行封装成一个匿名函数，再调用 cron 的 AddFunc 方法就可以了**
		* 将初始化的 Cron 对象放在根 Command 中。
		* 根 Command 结构中放入 Cron 实例，还放入了一个 CronSpecs 的数组，这个数组用来保存所有 Cron 命令的信息，为后续查看所有定时任务而准备
		* 在匿名函数中，封装的并不是传递进来的 Command，而是把这个 Command 做了一个副本，并且将其父节点设置为空，让它自身就是一个新的根节点；然后调用这个 Command 的 Execute 方法。从而不用根据参数进行遍历查询，提高效率。
```
// 创建一个cron实例
c := cron.New(cron.WithParser(cron.NewParser(cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)))

// 每秒执行一次
c.AddFunc("* * * * * *", func() { 
  fmt.Println("Every hour on the half hour") 
})
// 每秒调用一次Foo命令
rootCmd.AddCronCommand("* * * * * *", demo.FooCommand)
```
* 使用 cron 的三级命令对某个进程进行管理，要启动其他程序，创建进程
	* list        列出所有的定时任务  
	* restart     重启cron常驻进程  
	* start       启动cron常驻进程  
	* state       cron常驻进程状态  
	* stop        停止cron常驻进程
	* 可以使用标准库 osos.GetPid()获取pid
	* 在 Golang 中，如何启动一个当前二进制文件的子进程。
		* 使用开源**go-daemon**库，运行一个子进程，使用 os.StartProcess 来启动一个进程，执行当前进程相同的二进制文件以及当前进程相同的参数。调用一个 Reborn 方法启动一个子进程。
	* 启动三级命令行工具./fcou cron start -d，看到控制台打印出了子进程 PID 信息，和日志存储地址。
* 如何实现分布式定时器
	* appID 是为每个当前应用起的唯一标识，它用 Google 的 **uuid** 生成库 就可以很方便生成
	* 一个本地文件锁。当一个服务器上有多个进程需要进行抢锁操作，文件锁是一种单机多进程抢占的很简易的实现方式
	```
	多个进程同时使用 os.OpenFile 打开一个文件，并使用 syscall.Flock 带上 syscall.LOCK_EX 参数来对这个文件加文件锁，这里只会有一个进程抢占到文件锁，而其他抢占不到的进程从 syscall.Flock 函数中获取到的就是 error。根据这个 error 是否为空，我们就能判断是否抢占到了文件锁。
	```
	* 现在有了分布式选择器，就可以实现分布式调度了。我们为 Command 结构增加一个方法 AddDistributedCronCommand,它的实现和 AddCronCommand 差不多，唯一的区别就是在封装 cron.AddFunc 的匿名函数中，运行目标命令之前，要做一次分布式选举，如果被选举上了，才执行目标命令。
### 12｜配置和环境：配置服务中的设计思路
* 业务中有大量的配置项，比如数据库的用户名密码、缓存 Redis 的 IP 和端口、第三方调用的地址等。如何通过统一的方法快速获取到这些配置项。获取配置项的方法有很多，例如以下三种：
	* 读取本地配置文件
	* 读取远端配置服务
	* 获取环境变量
* **环境变量**，获取配置思路分析
	```
	环境变量（environment variables）一般是指在操作系统中用来指定操作系统运行环境的一些参数，如：临时文件夹位置和系统文件夹位置等。
	环境变量是在操作系统中一个具有特定名字的对象，它包含了一个或者多个应用程序所将使用到的信息。例如Windows和DOS操作系统中的path环境变量，当要求系统运行一个程序而没有告诉它程序所在的完整路径时，系统除了在当前目录下面寻找此程序外，还应到path中指定的路径去找。用户通过设置环境变量，来更好的运行进程。
	```
	* 在 baseFolder 目录下存放一个.env 文件保存默认值，这个.env 文件中的每一行，为一个环境变量 Key=Value 的形式，等运行的时候，进程运行的环境变量设置会覆盖.env 文件的设置。
		* 环境变量服务的接口及具体实现
		* 读取配置文件服务的接口及具体实现
	* 有的人可能会觉得要获取一个变量，直接使用配置文件就行啊，为什么要绕这么一圈？但是，环境变量是一个应用运行环境的一些参数，在现在的容器化流行的架构设计中，一般都会选择**使用环境变量来区别不同的环境和配置**。所以我们为框架提供获取环境变量的方式，在实际架构，特别是微服务相关的架构中，是非常有用的。
	* 环境变量中设置了 APP_ENV，表示这个应用在不同环境会有不同配置
	```
	APP_ENV=testing ./fcou env	
	```
* **配置**，多个配置文件按类别放在不同配置文件夹中
	* 配置服务的设计
		* 存储路径：最终的配置文件夹地址为，应用服务的 ConfigFolder 下的环境变量对应的文件夹，比如 ConfigFolder/development。
		* yaml格式
	* 配置文件的读取
		* **go-yaml** 就是非常通用的一个 Go 解析库
			* Marshal 表示序列化一个结构成为 YAML 格式；
			* Unmarshal 表示反序列化一个 YAML 格式文本成为一个结构；
			* 还有一个 UnmarshalStrict 函数，表示严格反序列化，比如如果 YAML 格式文件中包含重复 key 的字段，那么使用 UnmarshalStrict 函数反序列化会出现错误。
	* 配置文件的替换
		* 环境变量服务中，存放了包括.env 中设置的环境变量，那么我们自然会希望使用上这些环境变量，把配置文件中有的字段使用环境变量替换掉。
		* 那么这里在配置文件中就需要有一个“占位符”。这个占位符表示当前这个字段去环境变量中进行阅读。
	* 配置项的解析
	* 配置服务的代码实现
		* 先检查配置文件夹是否存在，然后读取文件夹中的每个以 yaml 或者 yml 后缀的文件；读取之后，先用 replace 对环境变量进行一次替换；替换之后使用 go-yaml，对文件进行解析。
	* 配置文件更新 App 服务
		* 在fcou框架的FcouApp实现中配置加载，configMap map[string]string
		* 在 configService，读取配置文件 loadConfigFile 的时候，加载 app.yaml 的配置文件的时候，就同时更新 FcouApp 里面的 configMap
	* 配置文件热更新
		* 我们可以自动监控配置文件目录下的所有文件，当配置文件有修改和更新的时候，能自动更新程序中的配置文件信息，也就是实现配置文件热更新。
		* 我们使用 **fsnotify库**能很方便对一个文件夹进行监控，当文件夹中有文件增 / 删 / 改的时候，会通过 channel 进行事件回调。
		* 这个库先使用 NewWatcher 创建一个监控器 watcher，然后使用 Add 来监控某个文件夹，通过 watcher 设置的 events 来判断文件是否有变化，如果有变化，就进行对应的操作，比如更新内存中配置服务存储的 map 结构。
*  Windows 系统不支持 pid 锁，所以我们需要在 Linux 或 Mac 系统下才能正常运行上面的程序。
----
### 13 日志：设计多输出的日志服务
* 日志等级设计，日志级别是按照从下往上的严重顺序排列的
	* panic，表示会导致整个程序出现崩溃的日志信息
	* fatal，表示会导致当前这个请求出现提前终止的错误信息
	* error，表示出现错误，但是不一定影响后续请求逻辑的错误信息
	* warn，表示出现错误，但是一定不影响后续请求逻辑的报警信息
	* info，表示正常的日志信息输出
	* debug，表示在调试状态下打印出来的日志信息
	* trace，表示最详细的信息，一般信息量比较大，可能包含调用堆栈等信息
* 日志格式
	* 日志级别，输出当前日志的级别信息。
	* 日志时间，输出当前日志的打印时间。
	* 日志简要信息，输出当前日志的简要描述信息，一句话说明日志错误。
	* 日志上下文字段，输出当前日志的附带信息。这些字段代表日志打印的上下文。
* 日志输出
	* 输出管道我们使用 io.Writer 来进行设置，接口设置
* 日志服务提供者
	* console.go 表示控制台输出，定义初始化实例方法 NewFcouConsoleLog；
	* single.go 表述单个日志文件输出，定义初始化实例方法 NewFcouSingleLog；
	* rotate.go 表示单个文件输出，但是自动进行切割，定义初始化实例方法 NewFcouRotateLog；
	* custom.go 表示自定义输出，定义实例化方法 NewFcouCustomLog。
* 日志服务的具体实现
	* 先判断日志级别是否符合要求，如果不符合要求，则直接返回，不进行打印；
	* 再使用 ctxFielder，从 context 中获取信息放在上下文字段中；
	* 接着将日志信息按照 formatter 序列化为字符串；
	* 最后通过 output 进行输出。
* 全链路日志
	* 在接收请求的时候，从请求 request 中解析全链路字段，存放进入 context 中
	* 在打印日志的时候从 context 中获取全链路字段序列化进入日志
	* 在发送请求的时候将全链路字段加入到 request 中
---
### 14:前后端一体化的架构方案
* 前后端分离的架构经常是，Nginx 作为网关，前端页面作为 Nginx 的一个 location 路由，而后端接口作为 Nginx 的另外一个 location 路由
	* 优点: 模块化，前端、后端、网关是独立模块，互相不干扰。
	* 缺点: 复杂度高一些，而且一旦在网关层增加逻辑，因为语言异构，网关层逻辑和业务层逻辑不是一个，开发起来会十分痛苦。
* Golang 的时代，我们是有另外一种选择的，那就是将 API 层往上提，使用 Golang 替代网关的逻辑
	* 优点：由于一体化了，开发复杂度降低了不少，开发效率提高，通用网关逻辑的研发成本就降低了
	* 缺点: 将三个模块一体化
* 前端Vue框架
	* Vue框架：“创造了一个概念叫渐进的框架。因为 Vue 的核心组成只是数据绑定和组件”
	* HTTP 库提供 FileServer 来封装对文件读取的 HTTP 服务。我们就可以使用这个 FileServer，来对前端编译出来的最终的浏览器可运行文件，提供 HTTP 服务
	* 将 Vue 的项目代码集成到业务代码中，然后确定其编译结果文件夹，在路由中，将某个请求路由到我们的编译结果文件夹，就完成了上面架构提到的同一个项目同时支持前端请求和后端请求了
	```
	路由注册时，先添加静态路由，再添加动态路由
	// /路径先去./dist目录下查找文件是否存在，找到使用文件服务提供服务
	r.Use(static.Serve("/", static.LocalFile("./fcou/dist", false)))
	```
	* 完整流程
		1. 安装npm,到node.js官网下载安装node.js,自带npm
		```
		node -v  //查看版本
		npm -v
		```
		2. 安装vue
		```
		npm install vue -g
		npm install vue-cli -g
		```
		3. 安装webpack
		```
		npm install  webpack -g
		npm install  webpack-cli -g
		```
		4. 新建项目，创建项目目录，进入此目录的终端
		```
		vue init webpack my-project   // 创建一个基于 webpack 模板的新项目
		// 这里需要进行一些配置，默认回车即可
		```
		5. 进入项目和安装依赖,生成的 index 文件，存放在根目录的 dist 目录下
		```
		npm install   //加载目录所需要的第三方库
		npm run build //执行 Vue 的编译操作
		```
		6. 启动项目
		```
		npm run dev
		```
	* vue-init webpack fcou 创建一个最完整的 fcou 项目
		* build 目录，存放项目构建（Webpack）相关的代码。
		* config 是配置目录，包括端口号等配置。
		* src 目录，这个是我们要开发的目录，业务逻辑代码基本上都在里面。
			* assets 存放一些图片信息
			* componenets 存放组件信息文件
			* App.vue 存放项目的入口组件 App
			* main.js 是项目的入口 js，引用加载入口组件 App。
		* static 目录，存放静态文件信息，比如图片、字体等。
		* test 目录，存放测试相关的信息。
		* dist 目录，编译完成后的结果
		* index.html 是 Vue 首页的入口页面
		* package.json 是项目的配置文件，保存引用的第三方库等信息
* 前后端一体化编译命令改造
	* 需求整理成以下四个命令
	```
	编译前端  ./fcou build frontend
	编译后端  ./fcou build backend
	同时编译前后端 ./fcou build all
	自编译 ./fcou build self
	```
---
### 15 提效：实现调试模式加速开发效率
* 后端开启调试模式
```
在使用 Vue 的时候， npm run dev 这个命令，为前端开启调试模式，只要你修改了 src 下的文件，编译器就会自动重新编译，并且更新浏览器页面上的渲染。那这种调试模式能否应用到 Golang 后端，进一步提升我们开发应用的效率。
```
	* 先通过 go build 编译出二进制文件，通过运行二进制文件再启动服务。
	* 对配置目录下的所有文件进行变更监听，一旦修改了配置目录下的文件，就重新更新内存中的配置文件 map
	* 前后端服务同时调试？是否能在前端和后端服务的前面，设计一个反向代理 proxy 服务呢
	* 让所有外部请求进入这个反响代理服务，然后由反向代理服务进行代理分发，前端请求分发到前端进程，后端请求分发到后端进程。
	* 一个请求到了，直接先请求一下后端服务，如果后端发现请求不存在，返回 404 Not Found 之后， 我们再将请求再请求到前端服务
	* 前端服务是直接使用 npm run dev 命令启动调试模式的
	* 后端服务是先进行  go build(监听文件有变化再编译) 再进行  go run ，所以后端服务是需要进行编译的
* 如何实现反向代理
	* 所谓反向代理，就是能将一个请求按照条件分发到不同的服务中去。在 Golang 中的 net/http/httputil 包中提供了 ReverseProxy 这么一个数据结构，它是实现整个反向代理的关键。
	```
	// 反向代理
	type ReverseProxy struct {
		// Director这个函数传入的参数是新复制的一个请求，我们可以修改这个请求
		// 比如修改请求的请求Host或者请求URL等
	Director func(*http.Request)

	// Transport 代表底层的连接池设置，比如连接最长保持多久等
		// 如果不填的话，则使用默认的设置
	Transport http.RoundTripper

	// FlushInterval表示多久将下游的response的数据拷贝到proxy的response
	FlushInterval time.Duration

	// ErrorLog 表示错误日志打印的句柄
	ErrorLog *log.Logger

	// BufferPool表示将下游response拷贝到proxy的response的时候使用的缓冲池大小
	BufferPool BufferPool

	// ModifyResponse 函数表示，如果要将下游的response内容进行修改，再传递给proxy
		// 的response，这个函数就可以进行设置，但是如果这个函数返回了error，则将response
		// 传递进入ErrorHandler，否则使用默认设置
	ModifyResponse func(*http.Response) error

	// ErrorHandler 处理ModifyResponse返回的Error
	ErrorHandler func(http.ResponseWriter, *http.Request, error)
	}
	```
	* Director 的参数是请求，表示如何对请求进行转发
	* ModifyResponse 字段，对下游的返回数据进行修改
	* ErrorHandler 有三个参数，responseWriter 是新 proxy 的 reponse 的写句柄，request 是 Director 修改后给下游的 request，而 error 则是 ModifyResponse 处理后的 error
* 配置项的设计
	* 我们先定义好默认的配置，然后从容器中获取配置服务，通过配置服务，获取对应的配置文件的设置，如果配置文件有对应字段的话，就进行对应字段的配置。
* 监控某个文件夹的变动，并且重新编译并且运行后端服务, 使用一种计时时间机制
	* 目标目录变更事件，有事件更新计时机制；
	* 计时机制到点事件，计时到点事件触发，代表有一个或多个目标目录变更已经存在，更新后端服务。
	```
	// 开启计时时间机制
	refreshTime := p.devConfig.Backend.RefreshTime
	t := time.NewTimer(time.Duration(refreshTime) * time.Second)
	// 先停止计时器
	t.Stop()
	for {
		select {
		case <-t.C:
			// 计时器时间到了，代表之前有文件更新事件重置过计时器
			// 即有文件更新
			fmt.Println("...检测到文件更新，重启服务开始...")
			if err := p.rebuildBackend(); err != nil {
				fmt.Println("重新编译失败：", err.Error())
			} else {
				if err := p.restartBackend(); err != nil {
					fmt.Println("重新启动失败：", err.Error())
				}
			}
			fmt.Println("...检测到文件更新，重启服务结束...")
			// 停止计时器
			t.Stop()
		case _, ok := <-watcher.Events:
			if !ok {
				continue
			}
			// 有文件更新事件，重置计时器
			t.Reset(time.Duration(refreshTime) * time.Second)
		case err, ok := <-watcher.Errors:
			if !ok {
				continue
			}
			// 如果有文件监听错误，则停止计时器
			fmt.Println("监听文件夹错误：", err.Error())
			t.Reset(time.Duration(refreshTime) * time.Second)
		}
	}
	```
---
### 16 自动化：自动化一切重复性劳动
* 自动化创建服务工具
	* 命令创建
		1. ./hade provider 一级命令，provider，打印帮助信息；
		2. ./hade provider new 二级命令，创建一个服务；
			* 借助一个第三方库 survey,实现命令行交互的方式
			* 已注册服务的字符串凭证比较，相同报错
			* 创建目录
			* 在目录下创建 contract.go、provider.go、service.go 三个文件
			* 利用text/template 库，在三个文件中根据预先定义好的模版填充内容
			```
			// 创建title这个模版方法
			funcs := template.FuncMap{"title": strings.Title}
			{
				//  创建contract.go
				file := filepath.Join(pFolder, folder, "contract.go")
				f, err := os.Create(file)
				if err != nil {
					return errors.Cause(err)
				}
				// 使用contractTmp模版来初始化template，并且让这个模版支持title方法，即支持{{.|title}}
				t := template.Must(template.New("contract").Funcs(funcs).Parse(contractTmp))
				// 将name传递进入到template中渲染，并且输出到contract.go 中
				if err := t.Execute(f, name); err != nil {
					return errors.Cause(err)
				}
			}
			template.New() 方法，创建一个 text/template 的 Template 结构，其中的参数 contract 字符串是为这个 Template 结构命名的，后面的 Funcs() 方法是将签名定义的模版函数注册到这个 Template 结构中，最后的 Parse() 是使用这个 Template 结构解析具体的模版文本。
			```
		3. ./hade provider list 二级命令，列出容器内的所有服务，列出它们的字符串凭证。
* 自动化创建命令行工具
	* 命令创建
		1. ./hade command 一级命令，显示帮助信息
		2. ./hade command list 二级命令，列出所有控制台命令
		3. ./hade command new 二级命令，创建一个控制台命令
* 自动化中间件迁移工具
	* 命令创建
		1. ./hade middleware 一级命令，显示帮助信息
		2. ./hade middleware list 二级命令，列出所有的业务中间件
		3. ./hade middleware new 二级命令，创建一个新的业务中间件
		4. ./hade middleware migrate 二级命令，迁移 Gin 已有的中间件
			* 参数中获取中间件名称；
			* 使用 go-git，将对应的 gin-contrib 的项目 clone 到目录 /app/http/middleware；
				* 从 git 上复制一个项目，在 Golang 中可以使用一个第三方库 go-git
				```
				_, err := git.PlainClone("/tmp/foo", false, &git.CloneOptions{
					URL:      "https://github.com/go-git/go-git",
					Progress: os.Stdout,
				})
				```
			* 删除不必要的文件 go.mod、go.sum、.git；
			* 替换关键字 “github.com/gin-gonic/gin”。
* 自动化初始化脚手架设计
	*  ./hade new 命令，要做的流程
		1. 下载 github.com/gohade/hade 的某个 release 版本到目标文件夹
			* go-github。这个库封装了 GitHub 的调用接口。
			```
			client := github.NewClient(nil)
			release, _, err = client.Repositories.GetLatestRelease(context.Background(), "gohade", "hade")
			```
			* 在返回的 RepositoryRelease 结构中，我们可以找到下载这个 release 版本的各种信息。其中包括 release 版本对应的版本号信息和 zip 下载地址。
			* 对于下载 zip 包，直接使用 http.Get 就能下载
			* 使用 Golang 标准库的 archive/zip，来读取 zip 包中的内容，然后将每个文件都复制到目标目录中
		2. 删除 framework 目录
		3. 修改 go.mod 中的模块名称
		4. 修改 go.mod 中的 require 信息，增加 require github.com/gohade/hade
		5. 修改所有文件使用业务目录的地方，将原本使用“github.com/Fcou/web-frame/app”  的所有引用改成 “[模块名称]/app”
	* ./hade new 命令，用户目前输入的三个信息：
		1. 目录名，最终是“当前执行目录 + 目录名”
		2. 模块名，最终创建应用的 module
		3. 版本号，对应的 hade 的 release 版本号