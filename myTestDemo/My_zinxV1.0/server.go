package main

import (
	"fmt"

	"github.com/hjy497373150/My_zinx/ziface"
	"github.com/hjy497373150/My_zinx/znet"
)

/*
	基于Zinx框架开发的 服务器端应用程序
*/

type PingRouter struct {
	znet.BaseRouter
} 

// test handle
func(this *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call PingRouter handle...")
	// 先读取客户端的数据，再写回ping..ping..ping..
	fmt.Println("Recv from client: msgId = ",request.GetMsgId(),",data = ",string(request.GetData()))

	err := request.GetConnection().SendMsg(100,[]byte("ping..ping..ping..."))
	if err != nil {
		fmt.Println(err)
	}
}

type HelloRouter struct {
	znet.BaseRouter
} 

// test handle
func(this *HelloRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call HelloRouter handle...")
	// 先读取客户端的数据，再写回ping..ping..ping..
	fmt.Println("Recv from client: msgId = ",request.GetMsgId(),",data = ",string(request.GetData()))

	err := request.GetConnection().SendMsg(200,[]byte("welcome to myzinx"))
	if err != nil {
		fmt.Println(err)
	}
}

// 创建链接后需要执行的hook函数
func DoConnBegin(conn ziface.IConnection) {
	fmt.Println("------> DoConnBegin is called")
	if err:=conn.SendMsg(202,[]byte("DoConnBegin!!!!!"));err!=nil {
		fmt.Println(err)
	}
	// 给当前的链接设置一些属性
	fmt.Println("Set conn Name,Github...")
	conn.SetProperty("Name","klayhu")
	conn.SetProperty("GitHub","https://github.com/hjy497373150/My_zinx")
	conn.SetProperty("E-mail","497373150@qq.com")
}

// 销毁链接前需要执行的hook函数
func DoConnLost(conn ziface.IConnection) {
	fmt.Println("------> DoConnLost is called")
	fmt.Println("connId = ",conn.GetConnID(),"is lost")

	// 获取链接属性
	if name, err := conn.GetProperty("Name");err==nil {
		fmt.Println("Name = ",name)
	}
	if github, err := conn.GetProperty("GitHub");err==nil {
		fmt.Println("GitHub = ",github)
	}
	if mail, err := conn.GetProperty("E-mail");err==nil {
		fmt.Println("E-mail = ",mail)
	}
}

func main() {
	// 1. 创建一个Server句柄，使用zinx的api
	s := znet.NewServer("[myzinx V1.0]")

	// 2.注册链接hook函数
	s.SetOnConnStart(DoConnBegin)
	s.SetOnConnStop(DoConnLost)

	// 3.给当前zinx框架添加一个自定义的Router
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloRouter{})

	// 4. 启动Server
	s.Serve()
}