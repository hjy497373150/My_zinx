package znet

import "github.com/hjy497373150/My_zinx/ziface"

type Request struct {
	Conn ziface.IConnection //已经和客户端建立好的链接
	Msg ziface.IMessage // 客户端请求的数据
}

// 得到当前的链接
func (r *Request)GetConnection() ziface.IConnection {
	return r.Conn
}

// 得到请求的消息数据
func (r *Request)GetData() []byte {
	return r.Msg.GetData()
}

// 得到消息ID
func (r *Request)GetMsgId() uint32 {
	return r.Msg.GetMsgId()
}