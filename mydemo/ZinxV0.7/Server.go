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

func main() {
	// 1、创建一个server 句柄， 使用Zinx的api
	s := znet.NewServer("[zinx v0.7]")

	// 2、给当前zinx框架添加一个router
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &PingRouter{})
	// 3、启动server
	s.Serve()
}
