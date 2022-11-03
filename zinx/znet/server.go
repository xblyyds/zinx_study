package znet

import (
	"fmt"
	"go-study/zinx/utils"
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
	// 当前 server 的消息模块管理, 用来绑定 MsgId 和对应的处理业务 API 关系
	MsgHandler ziface.IMsgHandle

	// 连接管理器
	ConnMgr ziface.IConnManager

	// 创建连接后自动调用的 Hook函数 OnConnStart
	OnConnStart func(conn ziface.IConnection)
	// 销毁连接后自动调用的 Hook函数 OnConnStop
	OnConnStop func(conn ziface.IConnection)
}

func (s *Server) Start() {
	fmt.Printf("[Zinx] Server Name: %s, listenner at IP: %s, Port: %d is starting",
		utils.GlobalObject.Name, utils.GlobalObject.Host, utils.GlobalObject.TcpPort)
	fmt.Printf("[Zinx] Version: %s, MaxConn: %d, MaxPackageSize: %d\n",
		utils.GlobalObject.Version, utils.GlobalObject.MaxConn, utils.GlobalObject.MaxPackageSize)

	fmt.Printf("[Start] Server Listener at IP :%s, Port :%d, is starting\n", s.IP, s.Port)
	go func() {
		// 0 开启消息队列及 worker 工作池
		s.MsgHandler.StartWorkerPool()

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
			// 3.1、阻塞等待客户端建立连接请求
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err ", err)
				continue
			}

			// 3.2、 设置服务器最大连接控制, 如果超过最大连接，那么则关闭此新的连接
			fmt.Println(">>>>>ConnMgr Len = ", s.ConnMgr.Len(), ", MAX = ", utils.GlobalObject.MaxConn)
			if s.ConnMgr.Len() >= utils.GlobalObject.MaxConn {
				// TODO 给客户端响应一个超出最大连接的错误包
				fmt.Println("====> Too Many Connections MaxConn = ", utils.GlobalObject.MaxConn)
				conn.Close()
				continue
			}

			// 3.3 TODO Server.Start() 处理该连接请求的 业务 方法, 此时应该有 handler 和 conn是绑定的

			// 将处理新链接的业务方法和conn进行绑定 得到我们的链接模块
			dealConn := NewConnection(s, conn, cid, s.MsgHandler)
			cid++

			// 启动当前的链接业务处理
			go dealConn.Start()
		}
	}()
}

func (s *Server) Stop() {
	fmt.Println("[STOP] Zinx server, name ", s.Name)
	// 将一些服务器的资源，状态或者一些已经开辟的链接信息，进行停止或者回收
	s.ConnMgr.ClearConn()
}

func (s *Server) Serve() {
	s.Start()

	//TODO 做一些启动服务器之后的额外业务

	select {}
}

func (s *Server) AddRouter(msgId uint32, router ziface.IRouter) {
	s.MsgHandler.AddRouter(msgId, router)
	fmt.Println("Add Router Succ!!")
}

func NewServer(name string) *Server {
	s := &Server{
		Name:       utils.GlobalObject.Name,
		IPVersion:  "tcp4",
		IP:         utils.GlobalObject.Host,
		Port:       utils.GlobalObject.TcpPort,
		MsgHandler: NewMsgHandle(),
		ConnMgr:    NewConnManager(),
	}
	return s
}

// 获取当前server 的连接管理器
func (s *Server) GetConnMgr() ziface.IConnManager {
	return s.ConnMgr
}

// 注册 OnConnStart 钩子函数
func (s *Server) SetOnConnStart(hookFunc func(conn ziface.IConnection)) {
	s.OnConnStart = hookFunc
}

// 调用 CallOnStart 钩子函数
func (s *Server) CallOnConnStart(conn ziface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("-----> Call OnConnStart()...")
		s.OnConnStart(conn)
	}
}

// 注册 OnConnStop 钩子函数
func (s *Server) SetOnConnStop(hookFunc func(conn ziface.IConnection)) {
	s.OnConnStop = hookFunc
}

// 调用 CallOnConnStop 函数
func (s *Server) CallOnConnStop(conn ziface.IConnection) {
	if s.OnConnStop != nil {
		fmt.Println("-----> Call OnConnStop...")
		s.OnConnStop(conn)
	}
}
