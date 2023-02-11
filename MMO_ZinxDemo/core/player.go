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

// 玩家广播世界聊天消息
func (p *Player) Talk(content string) {
	// 组建msgid = 200 的proto数据
	protoMsg := &pb.BroadCast{
		Pid: p.PlayerId,
		Tp: 1,
		Data: &pb.BroadCast_Content{
			Content: content,
		},
	}
	// 得到当前世界在线的所有玩家
	players := WorldMgrObj.GetAllPlayers()

	// 向所有玩家（包括自己）发送广播消息
	for _,player := range players {
		player.SendMsg(200,protoMsg)
	}
}

//同步玩家上线的位置消息
func (p *Player)SyncSurroundings() {
	// 1.获取当前玩家周围的玩家有哪些
	players := p.GetSurroundPlayers()

	// 2.将当前玩家的位置通过msgId = 200 发给周围的玩家（让其他玩家看到自己）
	// 2.1 组建msgId：200 的proto数据 
	protoMsg := &pb.BroadCast{
		Pid: p.PlayerId,
		Tp: 2,
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}
	// 2.2 全部周围的玩家都向格子的客户端发送200消息，protomsg
	for _,player := range players {
		player.SendMsg(200,protoMsg)
	}

	// 3.将周围的全部玩家的位置信息发送给当前玩家msgid = 202 客户端（让自己看到其他玩家）
	// 3.1组建MsgId:202 proto数据
		// 3.1.1 制作pb.player slice
	players_proto_msg := make([]*pb.Player, 0, len(players))
	for _,player := range players {
		// 制作一个message Player
		p:= &pb.Player{
			Pid: player.PlayerId,
			P: &pb.Position{
				X: player.X,
				Y: player.Y,
				Z: player.Z,
				V: player.V,
			},
		}
		players_proto_msg = append(players_proto_msg, p)
	}
		// 3.1.2 封装SyncPlayer protobuf数据
	SyncPlayer_proto_msg := &pb.SyncPlayers{
		Ps: players_proto_msg[:],
	}

	// 3.2将组建好的protomsg数据发给当前玩家的客户端
	p.SendMsg(202,SyncPlayer_proto_msg)

}

//广播当前玩家的位置移动信息
func (p *Player)UpdataPos(x,y,z,v float32) {
	// 1.更新当前玩家的坐标
	p.X = x
	p.Y = y
	p.Z = z
	p.V = v

	// 2.组件广播proto协议，msgid = 4 tp=4
	protoMsg := &pb.BroadCast{
		Pid: p.PlayerId,
		Tp: 4,
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}

	// 3.获得该玩家周围AOI九宫格之内的玩家
	players := p.GetSurroundPlayers()

	// 4.给周围的玩家发送当前玩家更新的消息
	for _,player := range players {
		player.SendMsg(200,protoMsg)
	}

}

// 获取当前玩家周边AOI九宫格之内的玩家
func (p *Player) GetSurroundPlayers() []*Player{
	// 得到当前玩家AOI九宫格之内的玩家id
	playerIds := WorldMgrObj.AoiManager.GetSurroundGridsByPos(p.X,p.Z)
	//  根据玩家id获取玩家
	players := make([]*Player,0,len(playerIds))
	for _,playerId := range playerIds {
		players = append(players, WorldMgrObj.GetPlayerByPid(int32(playerId)))
	}
	return players
}

func (p *Player) Offline() {
	// 1.得到当前玩家周边的九宫格中都有那些玩家
	players := p.GetSurroundPlayers()
	// 2.给周边玩家广播msgId = 201的消息
	protoMsg := &pb.SyncPid{
		Pid: p.PlayerId,
	}

	for _,player := range players {
		player.SendMsg(201,protoMsg)
	}
	// 3.将当前玩家从AOI管理器中删除
	WorldMgrObj.AoiManager.RemovePlayFromGridByPos(int(p.PlayerId),p.X,p.Z)
	// 4.将当前玩家从世界管理器删除
	WorldMgrObj.RemovePlayer(p.PlayerId)
}