package main

import (
	"fmt"

	"github.com/hjy497373150/My_zinx/MMO_ZinxDemo/apis"
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

	// 将当前新上线的玩家添加到WorldManager中
	core.WorldMgrObj.AddPlayer(player)

	// 将conn绑定一个playerid
	conn.SetProperty("playerId",player.PlayerId)

	// 同步周边玩家，告知他们当前玩家已上线，同步广播当前玩家的位置信息
	player.SyncSurroundings()


	fmt.Println("-------> Player Id = ",player.PlayerId ," 已上线")

}

// 当前客户端断开连接前的hook函数
func OnConnectionLost(conn ziface.IConnection) {
	// 1.通过conn得到当前链接的player
	playerId,_ := conn.GetProperty("playerId")
	player := core.WorldMgrObj.GetPlayerByPid(playerId.(int32))

	// 2.调用玩家下线业务
	player.Offline()
	

	fmt.Println("-------> Player Id = ",player.PlayerId ," 已下线")
}

// 当前客户端
func main() {
	// 创建zinx Server句柄
	s := znet.NewServer("MMO Game Based on Zinx")

	// 注册客户端建立连接和丢失函数
	s.SetOnConnStart(OnConnectionAdd)
	s.SetOnConnStop(OnConnectionLost)
	
	// 注册路由
	// 世界聊天
	s.AddRouter(2,&apis.WorldChatApi{})
	// 玩家移动
	s.AddRouter(3,&apis.Move{})
	
	//启动服务
	s.Serve()

	
}