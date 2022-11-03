package ziface

type IServer interface {
	// 启动服务器
	Start()
	// 停止服务器
	Stop()
	// 运行服务器
	Serve()
	// 路由功能，给当前的服务注册一个路由方法，供客户端的链接处理使用
	AddRouter(msgId uint32, router IRouter)

	// 获取当前server 的连接管理器
	GetConnMgr() IConnManager

	// 注册 OnConnStart 钩子函数
	SetOnConnStart(func(conn IConnection))

	// 调用 CallOnStart 钩子函数
	CallOnConnStart(conn IConnection)

	// 注册 OnConnStop 钩子函数
	SetOnConnStop(func(conn IConnection))

	// 调用 CallOnConnStop 函数
	CallOnConnStop(conn IConnection)
}
