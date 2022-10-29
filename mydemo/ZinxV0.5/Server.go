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
	if err := request.GetConnection().SendMsg(1, []byte("ping...ping...ping")); err != nil {
		fmt.Println(err)
	}
}

func main() {
	// 1、创建一个server 句柄， 使用Zinx的api
	s := znet.NewServer("[zinx v0.5]")

	// 2、给当前zinx框架添加一个router
	s.AddRouter(&PingRouter{})
	// 3、启动server
	s.Serve()
}
