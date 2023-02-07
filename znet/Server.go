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
	ConnMgr ziface.IConnManager // 当前Server的连接管理器
	
	// 新增两个hook函数
	// 当前Server 链接创建后调用的hook函数
	OnConnStart func(conn ziface.IConnection)
	// 当前Server 链接销毁之前调用的hook函数
	OnConnStop func(conn ziface.IConnection)
	
}


// 启动服务器
func (s *Server) Start() {
	fmt.Printf("[START] Server name: %s,listenner at IP: %s, Port %d is starting\n", s.Name, s.IP, s.Port)
	fmt.Printf("[Zinx] Version: %s, MaxConn: %d,  MaxPacketSize: %d, WorkPoolSize: %d\n",
		utils.GlobalObject.Version,
		utils.GlobalObject.MaxConn,
		utils.GlobalObject.MaxPackageSize,
		utils.GlobalObject.WorkerPoolSize)

	go func ()  {
		// 0.开启消息队列worker工作池
		s.MsgHandler.StartWorkPool()

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
		fmt.Println("Start Zinx server ", s.Name, "success, now listenning...")

		var cid uint32 = 0
		// 3.阻塞的等待客户端链接，处理客户端业务
		for {
			// 如果有客户端链接过来，阻塞会返回
			conn,err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Listener accept error ",err)
				continue
			}
			// 如果当前已有的链接超过了最大链接，就关闭新的链接
			if utils.GlobalObject.MaxConn <= s.ConnMgr.Len() {
				// TODO:给客户端响应一个超出最大链接的错误包
				fmt.Println("Too many conn,maxConn = ", utils.GlobalObject.MaxConn)
				conn.Close()
				continue
			} 

			dealconn := NewConnection(s, conn, cid, s.MsgHandler)
			cid++

			go dealconn.Start()
				
		}
	}()
	
}

// 停止服务器
func (s *Server) Stop() {
	// 将一些服务器的资源、状态或者一些已经开辟的链接信息 进行停止或回收
	fmt.Println("[STOP] Myzinx Server name:",s.Name)
	s.ConnMgr.ClearAll()
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

// 获取当前Server对应的connMgr
func (s *Server) GetConnMgr() ziface.IConnManager {
	return s.ConnMgr
}

// 注册该Server创建链接时的hook函数
func (s *Server) SetOnConnStart(hookFunc func (ziface.IConnection)) {
	s.OnConnStart = hookFunc
}

// 注册该Server销毁链接时的hook函数
func (s *Server) SetOnConnStop(hookFunc func (ziface.IConnection)) {
	s.OnConnStop = hookFunc
}

// 调用链接OnConnStart函数
func (s *Server) CallOnConnStart(conn ziface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("--->Call OnConnStart")
		s.OnConnStart(conn)
	}
}

// 调用链接OnConnStop函数
func (s *Server) CallOnConnStop(conn ziface.IConnection) {
	if s.OnConnStop != nil {
		fmt.Println("--->Call OnConnstop")
		s.OnConnStop(conn)
	}
}

func NewServer(name string) ziface.IServer {
	utils.GlobalObject.Reload() //加载配置模块
	s := &Server{
		Name:utils.GlobalObject.Name,
		IPVersion: "tcp4",
		IP: utils.GlobalObject.Host,
		Port: utils.GlobalObject.TcpPort,
		MsgHandler: NewMsgHandle(),
		ConnMgr: NewConnManager(),
	}
	return s
}