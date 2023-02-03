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
	fmt.Println("Call Router handle...")
	// 先读取客户端的数据，再写回ping..ping..ping..
	fmt.Println("Recv from client: msgId = ",request.GetMsgId(),",data = ",string(request.GetData()))

	err := request.GetConnection().SendMsg(1,[]byte("ping..ping..ping..."))
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	// 1. 创建一个Server句柄，使用zinx的api
	s := znet.NewServer("[myzinx V0.5]")
	// 2.给当前zinx框架添加一个自定义的Router
	s.AddRouter(new(PingRouter))
	// 3. 启动Server
	s.Serve()
}