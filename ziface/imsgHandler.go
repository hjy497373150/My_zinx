package ziface

/*
	消息管理抽象层
*/

type ImsgHandle interface {
	// 以非阻塞方式处理消息
	DoMsgHandle(request IRequest)
	// 为消息添加具体的处理逻辑
	AddRouter(msgId uint32, router IRouter)
	// 开启消息任务池
	StartWorkPool()
	// 将消息交给消息任务队列处理
	SendMsgToTaskQueue(request IRequest)
}