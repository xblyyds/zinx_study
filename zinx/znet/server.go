package znet

import (
	"fmt"
	"go-study/zinx/ziface"
	"net"
)

type Server struct {
	// 服务器的名称
	Name string
	// 服务器绑定的ip版本
	IPVersion string
	// 服务器监听的ip
	IP string
	// 服务器监听的端口
	Port int
	// 在当前的Server添加一个route, server注册的链接对应的处理业务
	Router ziface.IRouter
}

func (s *Server) Start() {
	fmt.Printf("[Start] Server Listener at IP :%s, Port :%d, is starting\n", s.IP, s.Port)
	go func() {
		// 1、获取一个TCP的Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr error : ", err)
			return
		}
		// 2、监听服务器的地址
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("list ", s.IPVersion, " err ", err)
			return
		}

		fmt.Println("start Zinx server succ, ", s.Name, " succ, Listening..")

		// 3、阻塞的等待客户端链接, 处理客户端链接业务(读写)
		var cid uint32
		cid = 0
		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err ", err)
				continue
			}

			// 将处理新链接的业务方法和conn进行绑定 得到我们的链接模块
			dealConn := NewConnection(conn, cid, s.Router)
			cid++

			// 启动当前的链接业务处理
			go dealConn.Start()
		}
	}()
}

func (s *Server) Stop() {
	// 将一些服务器的资源，状态或者一些已经开辟的链接信息，进行停止或者回收
}

func (s *Server) Serve() {
	s.Start()

	// 做一些启动服务器之后的额外业务

	select {}
}

func (s *Server) AddRouter(router ziface.IRouter) {
	s.Router = router
	fmt.Println("Add Router Succ!!")
}

func NewServer(name string) *Server {
	s := &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      8999,
		Router:    nil,
	}
	return s
}
