package main

import (
	"fmt"
	"go-study/mmo_game/api"
	"go-study/mmo_game/core"
	"go-study/zinx/ziface"
	"go-study/zinx/znet"
)

// 当客户端建立连接的时候的hook函数
func OnConnectionAdd(conn ziface.IConnection) {
	// 创建一个玩家
	player := core.NewPlayer(conn)

	// 同步当前的PlayerID给客户端, 走MsgID:1 消息
	player.SyncPID()

	// 同步当前玩家的初始化坐标信息给客户端, 走MsgID:200消息
	player.BroadCastStartPosition()

	// 新上线的玩家添加到 worldManager中
	core.WorldMgrObj.AddPlayer(player)
	//将该连接绑定属性PID
	conn.SetProperty("pID", player.PID)

	// 同步周边玩家上线信息
	player.SyncSurrounding()

	fmt.Println("=====> Player pID = ", player.PID, " arrived ====")
}

// 当客户端连接断开的时候的hook函数
func OnConnectionLost(conn ziface.IConnection) {
	// 获取当前连接的pID属性
	pID, _ := conn.GetProperty("pID")

	// 根据pID获取对应的玩家对象
	player := core.WorldMgrObj.GetPlayerByPID(pID.(int32))

	// 触发玩家下线业务
	if player != nil {
		player.LostConnection()
	}

	fmt.Println("====> Player ", pID, " left ====")
}

func main() {
	// 创建服务句柄
	s := znet.NewServer("Zinx MMO_GAME ")

	// 注册客户端连接简历 函数
	s.SetOnConnStart(OnConnectionAdd)
	s.SetOnConnStop(OnConnectionLost)

	//注册路由
	s.AddRouter(2, &api.WorldChatApi{}) // 世界聊天
	s.AddRouter(3, &api.MoveApi{})      // 玩家移动

	// 启动服务
	s.Serve()
}
