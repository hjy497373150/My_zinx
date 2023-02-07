package ziface

// 定义一个服务器接口
type IServer interface {
	// 启动服务器
	Start()
	// 停止服务器
	Stop()
	// 运行服务器
	Serve()

	// 路由功能：给当前的服务注册一个路由方法，供客户端链接处理使用
	AddRouter(msgId uint32, router IRouter)

	// 获取当前Server对应的connMgr
	GetConnMgr() IConnManager

	// 注册该Server创建链接时的hook函数
	SetOnConnStart(func (IConnection))
	// 注册该Server销毁链接时的hook函数
	SetOnConnStop(func (IConnection))
	// 调用链接OnConnStart函数
	CallOnConnStart(conn IConnection)
	// 调用链接OnConnStop函数
	CallOnConnStop(conn IConnection)
}