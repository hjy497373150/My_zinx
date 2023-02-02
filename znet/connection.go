package znet

import (
	"fmt"
	"net"

	"github.com/hjy497373150/My_zinx/ziface"
)

/*
	链接模块
*/

type Connection struct {
	Conn *net.TCPConn // 当前链接的TCP socket套接字
	ConnID uint32 // 链接的ID
	isClosed bool // 当前的链接状态（是否关闭
	ExitChan chan bool // 告知当前链接已经退出/停止 channel
	Router ziface.IRouter // 当前链接的处理方法
}

// 初始化链接模块的方法
func NewConnection(conn *net.TCPConn, connID uint32, router ziface.IRouter) *Connection {
	c := &Connection {
		Conn: conn,
		ConnID: connID,
		isClosed: false,
		ExitChan: make(chan bool, 1),
		Router: router,
	}

	return c
}

func (c *Connection)StartReader() {
	fmt.Println("Reader Gorountine is running...")
	defer fmt.Println("connID = ",c.ConnID, " Reader is exit, remote addr is ", c.RemoteAddr().String())
	defer c.Stop() // defer工作原理 多个defer按栈的原理先后执行
	for {
		// 读取客户端的数据到buf中
		buf := make([]byte, 1024)
		_, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("Receive buf error",err)
			continue
		}

		// 得到当前conn数据的Request请求数据
		req := Request{
			Conn: c,
			Data: buf,
		}

		// 从路由中，找到绑定的conn对应的router调用
		go func(request ziface.IRequest) {
			c.Router.PreHandle(request)
			c.Router.Handle(request)
			c.Router.PostHandle(request)
		}(&req)

	}
}

// 启动链接，让当前的链接准备开始工作
func (c *Connection)Start(){
	fmt.Println("Conn Start()... connID = ",c.ConnID)

	// 启动从当前链接读数据的业务
	go c.StartReader()
} 

// 关闭链接，结束当前链接工作
func (c *Connection)Stop() {
	fmt.Println("Conn Stop()... connID = ",c.ConnID)

	if c.isClosed == true {
		return //当前链接已经关闭直接返回
	}

	c.isClosed = true

	close(c.ExitChan)
}

// 获取当前绑定的TCP socket conn
func (c *Connection)GetTcpConnection() *net.TCPConn  {
	return c.Conn
}

// 获取当前链接模块的链接ID
func (c *Connection)GetConnID() uint32{
	return c.ConnID
}

// 获取远程客户端的TCP状态
func (c *Connection)RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// 发送数据给远程的客户端
func (c *Connection) Send(data []byte) error {
	return nil
}