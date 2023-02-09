package main

import (
	"fmt"

	"github.com/hjy497373150/My_zinx/MMO_ZinxDemo/core"
	"github.com/hjy497373150/My_zinx/ziface"
	"github.com/hjy497373150/My_zinx/znet"
)

/*
	服务器的主入口
*/

// 当前客户端建立连接后的hook函数

func OnConnectionAdd(conn ziface.IConnection) {
	// 创建一个player对象
	player := core.NewPlayer(conn)

	// 给客户端发送msgid = 1的消息，同步当前player的id给客户端
	player.SyncPid()

	// 给客户端发送msgid = 200的消息，同步当前player的初始位置给客户端
	player.BroadCastStartPosition()

	fmt.Println("-------> Player Id = ",player.PlayerId ," 已上线")
}

func main() {
	// 创建zinx Server句柄
	s := znet.NewServer("MMO Game Based on Zinx")

	// 注册客户端建立连接和丢失函数
	s.SetOnConnStart(OnConnectionAdd)
	
	//启动服务
	s.Serve()
}