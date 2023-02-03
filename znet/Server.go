package znet

import (
	"fmt"
	"net"

	"github.com/hjy497373150/My_zinx/utils"
	"github.com/hjy497373150/My_zinx/ziface"
)

// IServer 接口实现
type Server struct {
	Name string // 服务器名称
	IPVersion string //服务器绑定的ip版本
	IP string // 服务器监听的ip
	Port int // 服务器监听的端口
	MsgHandler ziface.ImsgHandle // 当前Server的消息管理模块，用来绑定MsgId和对应处理业务的Api关系
}


// 启动服务器
func (s *Server) Start() {
	fmt.Printf("[START] Server name: %s,listenner at IP: %s, Port %d is starting\n", s.Name, s.IP, s.Port)
	fmt.Printf("[Zinx] Version: %s, MaxConn: %d,  MaxPacketSize: %d\n",
		utils.GlobalObject.Version,
		utils.GlobalObject.MaxConn,
		utils.GlobalObject.MaxPackageSize)

	go func ()  {
		// 1.获取一个TCP的addr
		addr,err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d",s.IP,s.Port))
		if err != nil {
			fmt.Println("Resolve TCP Addr error :",err)
			return
		}

		// 2.监听服务器的地址
		listener,err := net.ListenTCP(s.IPVersion,addr)
		if err != nil {
			fmt.Println("Listen ",s.IPVersion,"error ",err)
			return
		}
		//已经监听成功
		fmt.Println("Start Zinx server  ", s.Name, " success, now listenning...")

		var cid uint32 = 0
		// 3.阻塞的等待客户端链接，处理客户端业务
		for {
			// 如果有客户端链接过来，阻塞会返回
			conn,err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Listener accept error ",err)
				continue
			}

			dealconn := NewConnection(conn, cid, s.MsgHandler)
			cid++

			go dealconn.Start()

		}
	}()
	
}

// 停止服务器
func (s *Server) Stop() {
	// TODO:将一些服务器的资源、状态或者一些已经开辟的链接信息 进行停止或回收
}

// 运行服务器
func (s *Server) Serve() {
	// 启动server的服务功能
	s.Start()

	// TODO:做一些启动服务器后的额外业务

	//阻塞状态
	select{}
}

// 路由功能：给当前的服务注册一个路由方法，供客户端链接处理使用
func (s *Server) AddRouter(msgId uint32, router ziface.IRouter) {
	s.MsgHandler.AddRouter(msgId,router)

	fmt.Println("Add Router success...")	
}

func NewServer(name string) ziface.IServer {
	utils.GlobalObject.Reload() //加载配置模块
	s := &Server{
		Name:utils.GlobalObject.Name,
		IPVersion: "tcp4",
		IP: utils.GlobalObject.Host,
		Port: utils.GlobalObject.TcpPort,
		MsgHandler: NewMsgHandle(),
	}
	return s
}