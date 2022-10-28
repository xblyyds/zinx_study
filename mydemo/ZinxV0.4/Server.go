package main

import (
	"fmt"
	"go-study/zinx/ziface"
	"go-study/zinx/znet"
)

type PingRouter struct {
	znet.BaseRouter
}

func (this *PingRouter) PreHandle(request ziface.IRequest) {
	fmt.Println("Call Router PreHandle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("before ping...\n"))
	if err != nil {
		fmt.Println("call back before ping error")
	}
}

func (this *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Router Handle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("ping...ping...ping\n"))
	if err != nil {
		fmt.Println("call back ping...ping...ping error")
	}
}

func (this *PingRouter) PostHandle(request ziface.IRequest) {
	fmt.Println("Call Router PostHandle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("after ping\n"))
	if err != nil {
		fmt.Println("call back after ping error")
	}
}

func main() {
	// 1、创建一个server 句柄， 使用Zinx的api
	s := znet.NewServer("[zinx v0.4]")

	// 2、给当前zinx框架添加一个router
	s.AddRouter(&PingRouter{})
	// 3、启动server
	s.Serve()
}
