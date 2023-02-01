package znet

import (
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
	handleAPI ziface.HandleFunc //当前链接所绑定的处理业务的方法API
	ExitChan chan bool // 告知当前链接已经退出/停止 channel
}

// 初始化链接模块的方法
func NewConnection(conn *net.TCPConn, connID uint32, callback_api ziface.HandleFunc) *Connection {
	c := &Connection {
		Conn: conn,
		ConnID: connID,
		isClosed: false,
		handleAPI: callback_api,
		ExitChan: make(chan bool, 1),
	}
	
	return c
}