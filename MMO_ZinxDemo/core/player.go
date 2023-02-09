package core

import (
	"fmt"
	"math/rand"
	"sync"
	"github.com/hjy497373150/My_zinx/MMO_ZinxDemo/pb"
	"github.com/hjy497373150/My_zinx/ziface"
	"google.golang.org/protobuf/proto"
)

type Player struct {
	PlayerId int32 // 玩家Id
	Conn	ziface.IConnection // 当前玩家与客户端之间的链接
	X float32 // 平面的X坐标
	Y float32 // 高度
	Z float32 // 平面的y坐标 注意不是Y
	V float32 // 旋转的角度 0-360°
}

/*
	PLayerID 生成器
*/
var PidGen int32 = 1 //用来生成玩家Id的计数器
var IdLock sync.Mutex //保护PidGen的锁

// 创建一个玩家的方法
func NewPlayer(conn ziface.IConnection) *Player {
	// 生成一个玩家Id
	IdLock.Lock()
	id := PidGen
	PidGen++
	IdLock.Unlock()

	// 创建一个玩家对象
	p := &Player{
		PlayerId: id,
		Conn: conn,
		X: float32(160 + rand.Intn(10)),
		Y: 0,
		Z: float32(140 + rand.Intn(20)),
		V: 0,
	}

	return p
}

/*
	提供一个发送给客户端消息的方法
	主要是将pb的protobuf数据序列化之后，再调用zinx的sendmsg方法
*/
func (p *Player) SendMsg(msgId uint32, data proto.Message) {
	// 将proto message结构体序列化 转换成二进制
	msg, err:=proto.Marshal(data)
	if err != nil {
		fmt.Println("Marshal msg error",err)
		return 
	}
	// 将二进制文件 通过zinx 的sendmsg将数据发送给客户端
	if p.Conn == nil {
		fmt.Println("conn in player is nil")
		return 
	}

	if err := p.Conn.SendMsg(msgId,msg);err != nil {
		fmt.Println("player sendmsg error ",err)
		return 
	}

}

// 告知客户端玩家playerid，同步已经生成的玩家id给客户端
func (p *Player) SyncPid() {
	// 组建msgid = 1 的proto数据
	protoMsg := &pb.SyncPid{
		Pid:p.PlayerId,
	}

	p.SendMsg(1,protoMsg)
}

// 广播玩家自己的出生地点
func (p *Player) BroadCastStartPosition() {
	// 组建msgid = 200 的proto数据
	protoMsg := &pb.BroadCast{
		Pid: p.PlayerId,
		Tp: 2, // 2代表广播的位置坐标
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}

	p.SendMsg(200,protoMsg)
}