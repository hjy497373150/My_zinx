package znet

import (
	"fmt"
	"net"

	"github.com/hjy497373150/My_zinx/ziface"
)

// IServer 接口实现
type Server struct {
	Name string // 服务器名称
	IPVersion string //服务器绑定的ip版本
	IP string // 服务器监听的ip
	Port int // 服务器监听的端口
	Router ziface.IRouter // 当前Server由用户绑定的回调router，server注册的链接对应的处理业务
}


// 启动服务器
func (s *Server) Start() {
	fmt.Printf("[Start] Server Listener at IP : %s, Port: %d is starting\n",s.IP, s.Port)

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

		var cid uint32
		cid = 0
		// 3.阻塞的等待客户端链接，处理客户端业务
		for {
			// 如果有客户端链接过来，阻塞会返回
			conn,err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Listener accept error ",err)
				continue
			}

			dealconn := NewConnection(conn, cid, s.Router)
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
func (s *Server) AddRouter(router ziface.IRouter) {
	s.Router = router
	fmt.Println("Add Router success...")

}

func NewServer(name string)  ziface.IServer {
	s := &Server{
		Name: name,
		IPVersion: "tcp4",
		IP: "127.0.0.1",
		Port: 8999,
		Router: nil,
	}
	return s
}