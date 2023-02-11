package apis

import (
	"fmt"

	"github.com/hjy497373150/My_zinx/MMO_ZinxDemo/core"
	"github.com/hjy497373150/My_zinx/MMO_ZinxDemo/pb"
	"github.com/hjy497373150/My_zinx/ziface"
	"github.com/hjy497373150/My_zinx/znet"
	"google.golang.org/protobuf/proto"
)

/*
	玩家移动
*/
type Move struct {
	znet.BaseRouter
}

func (mv *Move) Handle(request  ziface.IRequest) {
	// 1.解析客户端传递进来的proto协议
	protomsg := &pb.Position{}
	err := proto.Unmarshal(request.GetData(), protomsg)
	if err != nil {
		fmt.Println("Move:Position Unmarshal error",err)
		return
	}

	// 2.得到当前发送位置的是哪个玩家
	playerId,err := request.GetConnection().GetProperty("playerId")
	if err != nil {
		fmt.Println("GetProperty playerId error",err)
		return
	}

	fmt.Printf("playerId = %d, move (%f,%f,%f,%f)",playerId,protomsg.X,protomsg.Y,protomsg.Z,protomsg.V)

	// 3.给其他玩家进行当前玩家的位置信息广播
	player := core.WorldMgrObj.GetPlayerByPid(playerId.(int32))

	// 4. 广播并更新当前玩家的坐标
	player.UpdataPos(protomsg.X,protomsg.Y,protomsg.Z,protomsg.V)
}