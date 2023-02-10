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
	世界聊天 路由业务
*/	

type WorldChatApi struct {
	znet.BaseRouter
}

func (wc *WorldChatApi) Handle(request ziface.IRequest) {
	// 1.解析客户端传递进来的proto协议
	protomsg := &pb.Talk{}
	err := proto.Unmarshal(request.GetData(), protomsg)
	if err != nil {
		fmt.Println("Talk Unmarshal error",err)
		return
	}

	// 2.当前的聊天数据是属于哪个玩家发送的
	playerId,err := request.GetConnection().GetProperty("playerId")
	if err != nil {
		fmt.Println("GetProperty playerId error",err)
		request.GetConnection().Stop()
		return
	}

	// 3.根据playerId得到对应的player对象
	player := core.WorldMgrObj.GetPlayerByPid(playerId.(int32))

	// 将这个消息广播给其他全部在线的对象
	player.Talk(protomsg.Content)

}