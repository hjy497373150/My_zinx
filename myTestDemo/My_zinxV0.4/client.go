package main

import (
	"fmt"
	"net"
	"time"
)

/*
	模拟客户端
*/
func main() {
	fmt.Println("Client start...")

	time.Sleep(1 * time.Second)

	// 1.链接远程服务器，得到一个conn链接
	conn,err := net.Dial("tcp","127.0.0.1:8888")
	if err != nil {
		fmt.Println("Client start err, exit!")
		return
	}

	// 2.链接调用write写数据
	for {
		_,err := conn.Write([]byte("hello MyzinxV0.4..."))
		if err != nil {
			fmt.Println("Client write conn err",err)
			continue
		}

		// 3.接收服务器的回显
		buf := make([]byte, 1024)
		cnt,err := conn.Read(buf)
		if err != nil {
			fmt.Println("Client read buf error ",err)
			return
		}

		fmt.Printf("Server call back: %s ,cnt = %d\n",buf,cnt)

		// 阻塞一下，1s发一次
		time.Sleep(1 * time.Second)
	}
	

}