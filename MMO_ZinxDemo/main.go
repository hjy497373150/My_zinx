package main

import "github.com/hjy497373150/My_zinx/znet"

/*
	服务器的主入口
*/

func main() {
	// 创建zinx Server句柄
	s := znet.NewServer("MMO Game Based on Zinx")

	//启动服务
	s.Serve()
}