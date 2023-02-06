package znet

import (
	"errors"
	"fmt"
	"io"
	"net"

	"github.com/hjy497373150/My_zinx/utils"
	"github.com/hjy497373150/My_zinx/ziface"
)

/*
	链接模块
*/

type Connection struct {
	Conn *net.TCPConn // 当前链接的TCP socket套接字
	ConnID uint32 // 链接的ID
	isClosed bool // 当前的链接状态（是否关闭
	ExitChan chan bool // 告知当前链接已经退出/停止 channel
	MsgHandle ziface.ImsgHandle// 消息的管理msgId和对应的处理业务Api关系
	MsgChan chan []byte // 读写分离使用，用于读写gorountine通信的无缓冲数据通道
}

// 初始化链接模块的方法
func NewConnection(conn *net.TCPConn, connID uint32, msgHandle ziface.ImsgHandle) *Connection {
	c := &Connection {
		Conn: conn,
		ConnID: connID,
		isClosed: false,
		ExitChan: make(chan bool, 1),
		MsgHandle: msgHandle,
		MsgChan: make(chan []byte),
	}

	return c
}

func (c *Connection)StartReader() {
	fmt.Println("[Reader Gorountine is running...]")
	defer fmt.Println("[Reader is exit] connID = ",c.ConnID, ", remote addr is ", c.RemoteAddr().String())
	defer c.Stop() // defer工作原理 多个defer按栈的原理先后执行
	for {
		// 读取客户端的数据到buf中

		// 创建封包拆包的对象
		dp := NewDataPack()
		// 读取客户端Msg Head二进制流 8字节
		headData := make([]byte, dp.GetHeadLen())
		if _,err := io.ReadFull(c.Conn,headData);err!=nil {
			fmt.Println("Read headData error ",err)
			break
		} 

		// 拆包得到datalen和id 放在msg消息中
		msg,err := dp.UnPack(headData)
		if err != nil {
			fmt.Println("Unpack error...",err)
			break
		}

		// 根据datalen再次读取data 放在msg.data中
		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			if _,err := io.ReadFull(c.Conn,data);err!=nil {
				fmt.Println("Read msg data error ",err)
				break
			}
		}
		msg.SetData(data)
		// 得到当前conn数据的Request请求数据
		req := Request{
			Conn: c,
			Msg: msg,
		}
		
		
		if utils.GlobalObject.WorkerPoolSize > 0 {
			// 已经开启了工作池机制，将消息发给Worker工作池处理
			c.MsgHandle.SendMsgToTaskQueue(&req)
		} else {
			// 从绑定好的消息和对应的处理方法中执行对应的Handle方法
			go c.MsgHandle.DoMsgHandle(&req)
		}

	}
}

/*
	写消息goroutine，将数据从channel中拿出来发给客户端
*/
func (c *Connection)StartWriter() {
	fmt.Println("[Writer Gorountine is running...]")
	defer fmt.Println("[Writer is exit] connID = ",c.ConnID, ", remote addr is ", c.RemoteAddr().String())

	for {
		select {
		case data := <-c.MsgChan:
			// 有数据要发给客户端
			if _,err := c.Conn.Write(data);err != nil {
				fmt.Println("send msg error ",err," Conn Writer exit")
				return
			}
		case <- c.ExitChan:
			// Conn已经关闭（由StartReader结束的defer调用c.stop确定
			return
		}
	}
}

// 启动链接，让当前的链接准备开始工作
func (c *Connection)Start(){
	fmt.Println("Conn Start()... connID = ",c.ConnID)

	// 启动从当前链接读数据的业务
	go c.StartReader()

	// 启动从当前链接写数据的业务
	go c.StartWriter()
} 

// 关闭链接，结束当前链接工作
func (c *Connection)Stop() {
	fmt.Println("Conn Stop()... connID = ",c.ConnID)

	if c.isClosed == true {
		return //当前链接已经关闭直接返回
	}

	c.isClosed = true

	// 关闭socket链接
	c.Conn.Close()
	// 告知writer关闭
	c.ExitChan <- true
	// 回收资源
	close(c.ExitChan)
	close(c.MsgChan)
}

// 获取当前绑定的TCP socket conn
func (c *Connection)GetTcpConnection() *net.TCPConn  {
	return c.Conn
}

// 获取当前链接模块的链接ID
func (c *Connection)GetConnID() uint32{
	return c.ConnID
}

// 获取远程客户端的TCP状态
func (c *Connection)RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// 发送数据给远程的客户端
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("Conn is closed when send msg")
	}

	// 将data进行封包
	dp := NewDataPack()

	binaryMsg,err := dp.Pack(NewMessage(msgId,data))
	if err != nil {
		fmt.Println("Pack error msg id = ",msgId)
		return errors.New("Pack error msg")
	}


	// 将数据发给Msgchan
	c.MsgChan <- binaryMsg 


	return nil
}