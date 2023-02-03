package main

import (
	"fmt"
	"io"
	"net"
	"time"

	"github.com/hjy497373150/My_zinx/znet"
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
		// _,err := conn.Write([]byte("hello MyzinxV0.5..."))
		// if err != nil {
		// 	fmt.Println("Client write conn err",err)
		// 	continue
		// }

		// // 3.接收服务器的回显
		// buf := make([]byte, 1024)
		// cnt,err := conn.Read(buf)
		// if err != nil {
		// 	fmt.Println("Client read buf error ",err)
		// 	return
		// }

		// fmt.Printf("Server call back: %s ,cnt = %d\n",buf,cnt)

		// 发送封包的message消息 msgid = 0
		dp := znet.NewDataPack()
		binaryMsg,err := dp.Pack(znet.NewMessage(0,[]byte("MyZinxV0.5 client test message")))
		if err != nil {
			fmt.Println("Client pack msg error ",err)
			return
		}
		if _,err := conn.Write(binaryMsg);err!=nil {
			fmt.Println("Client write error")
			continue
		}

		// 服务器应该回复一个msg回显 ，msgid=1 pingpingping
		// 先读取流中的head部分 得到msgid和msglendata
		binaryHead := make([]byte,dp.GetHeadLen())
		if _,err := io.ReadFull(conn,binaryHead);err!=nil {
			fmt.Println("Client read head error ",err)
			break
		}

		msgHead,err:=dp.UnPack(binaryHead)
		if err != nil {
			fmt.Println("client unpack msgHead error ",err)
			break
		}

		// 根据datalen进行第二次读取，将数据读出来
		if msgHead.GetMsgLen() > 0 {
			msg := msgHead.(*znet.Message)
			msg.Data = make([]byte, msg.GetMsgLen())
			if _,err := io.ReadFull(conn,msg.Data);err!=nil {
				fmt.Println("Read msg data error ",err)
				return
			}

			fmt.Println("------->recv Server Msg: ID = ",msg.ID,", len = ",msg.DataLen,",data = ",string(msg.Data))
		}

		// 阻塞一下，1s发一次
		time.Sleep(1 * time.Second)

	}
	

}