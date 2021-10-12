# web-frame
#### 跟着轩脉刃从 0 开始构建 Web 框架
##### step_01
- Web Server 的本质
> 实际上就是接收、解析 HTTP 请求传输的文本字符，理解这些文本字符的指令，然后进行计算，再将返回值组织成 HTTP 响应的文本字符，通过 TCP 网络传输回去。
- 使用go语言的net/http标准库，阅读学习源码的方法：库函数 > 结构定义 > 结构函数。
- net/http库 Server 整个逻辑线大致是:
> 创建服务 -> 监听请求 -> 创建连接 -> 处理请求
> 第一层，标准库创建 HTTP 服务是通过创建一个 Server 数据结构完成的；
> 第二层，Server 数据结构在 for 循环中不断监听每一个连接；
> 第三层，每个连接默认开启一个 Goroutine 为其服务；
> 第四、五层，serverHandler 结构代表请求对应的处理逻辑，并且通过这个结构进行具体业务逻辑处理；
> 第六层，Server 数据结构如果没有设置处理函数 Handler，默认使用 DefaultServerMux 处理请求；
> 第七层，DefaultServerMux 是使用 map 结构来存储和查找路由规则。
- 创建 Server 数据结构，并且在数据结构中创建了自定义的 Handler（Core 数据结构）和监听地址，实现了一个 HTTP 服务。
