package ziface

import "net"

// 定义链接模块的抽象层
type IConnection interface {
	Start() // 启动链接，让当前的链接准备开始工作
	Stop() // 关闭链接，结束当前链接工作
	GetTcpConnection() *net.TCPConn // 获取当前绑定的TCP socket conn
	GetConnID() // 获取当前链接模块的链接ID
	RemoteAddr() net.Addr // 获取远程客户端的TCP状态
	Send(data []byte) error // 发送数据给远程的客户端
}

// 定义一个处理链接业务的方法
type HandleFunc func(*net.TCPConn, []byte, int) error