# web-frame
#### 跟着轩脉刃从 0 开始构建 Web 框架
---
### step_01 分析 net/http，创建 Server 数据结构，自实现Handler接口
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
### step_02 添加上下文 Context 为请求设置超时时间
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

