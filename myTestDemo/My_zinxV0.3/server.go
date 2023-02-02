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

// test prehandle
func(this *PingRouter) PreHandle(request ziface.IRequest) {
	fmt.Println("Call Router Prehandle...")
	_,err := request.GetConnection().GetTcpConnection().Write([]byte("before ping\n"))
	if err != nil {
		fmt.Println("Call back before ping ... error",err)
	}
}

// test handle
func(this *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Router handle...")
	_,err := request.GetConnection().GetTcpConnection().Write([]byte("ping ping\n"))
	if err != nil {
		fmt.Println("Call back ping ping ... error",err)
	}
}

// test posthandle
func(this *PingRouter) PostHandle(request ziface.IRequest) {
	fmt.Println("Call Router Posthandle...")
	_,err := request.GetConnection().GetTcpConnection().Write([]byte("after ping\n"))
	if err != nil {
		fmt.Println("Call back after ping ... error",err)
	}
}


func main() {
	// 1. 创建一个Server句柄，使用zinx的api
	s := znet.NewServer("[myzinx V0.3]")
	// 2.给当前zinx框架添加一个自定义的Router
	s.AddRouter(new(PingRouter))
	// 3. 启动Server
	s.Serve()
}