package main

import (
	"fmt"
	"go-study/zinx/ziface"
	"go-study/zinx/znet"
)

type PingRouter struct {
	znet.BaseRouter
}

func (this *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Router Handle...")
	// 先读取客户端的数据, 再回写 ping...ping...ping
	fmt.Println("recv from client: msgId = ", request.GetMsgId(), ", data = ", string(request.GetData()))
	if err := request.GetConnection().SendMsg(200, []byte("Hello Welcome to Zinx!")); err != nil {
		fmt.Println(err)
	}
}

type HelloZinxRouter struct {
	znet.BaseRouter
}

func (h *HelloZinxRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call HelloZinRouter Handle...")
	// 先读取客户端的数据, 再回写 ping...ping...ping
	fmt.Println("recv from client: msgId = ", request.GetMsgId(), ", data = ", string(request.GetData()))
	if err := request.GetConnection().SendMsg(201, []byte("ping...ping...ping")); err != nil {
		fmt.Println(err)
	}
}

// 创建连接之后执行的钩子函数
func DoConnectionBegin(conn ziface.IConnection) {
	fmt.Println("=====> DoConnectionBegin is Called")
	if err := conn.SendMsg(202, []byte("DoConnection BEGIN!")); err != nil {
		fmt.Println(err)
	}

	// 给当前连接设置一些属性
	fmt.Println("Set conn property...")
	conn.SetProperty("Name", "徐步亮")
	conn.SetProperty("GitHub", "https://github.com/xblyyds")
}

// 连接断开的钩子函数
func DoConnectionStop(conn ziface.IConnection) {
	fmt.Println("=====> DoConnectionStop is Called")
	fmt.Println("connID = ", conn.GetConnID(), " is Stop...")

	// 获取连接属性
	if name, err := conn.GetProperty("Name"); err == nil {
		fmt.Println("Name = ", name)
	}
	if github, err := conn.GetProperty("GitHub"); err == nil {
		fmt.Println("GitHub = ", github)
	}
}

func main() {
	// 1、创建一个server 句柄， 使用Zinx的api
	s := znet.NewServer("[zinx v0.9]")

	// 2、注册连接 Hook钩子函数
	s.SetOnConnStart(DoConnectionBegin)
	s.SetOnConnStop(DoConnectionStop)

	// 3、给当前zinx框架添加一个router
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloZinxRouter{})
	// 4、启动server
	s.Serve()
}
