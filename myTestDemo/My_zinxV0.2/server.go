package main

import (
	"github.com/hjy497373150/My_zinx/znet"
)

/*
	基于Zinx框架开发的 服务器端应用程序
*/

func main() {
	// 1. 创建一个Server句柄，使用zinx的api
	s := znet.NewServer("[myzinx V0.2]")
	// 2. 启动Server
	s.Serve()
}