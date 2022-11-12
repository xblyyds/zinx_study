package api

import (
	"fmt"
	"go-study/mmo_game/core"
	"go-study/mmo_game/pb"
	"go-study/zinx/ziface"
	"go-study/zinx/znet"
	"google.golang.org/protobuf/proto"
)

// 玩家移动
type MoveApi struct {
	znet.BaseRouter
}

func (m *MoveApi) Handle(request ziface.IRequest) {
	// 1. 将客户端传来的proto协议编码
	msg := &pb.Position{}
	err := proto.Unmarshal(request.GetData(), msg)
	if err != nil {
		fmt.Println("Move: Position Unmarshal error ", err)
		return
	}
	// 2. 得知当前的消息是从哪个玩家传递来的,从连接属性pID中获取
	pID, err := request.GetConnection().GetProperty("pID")
	if err != nil {
		fmt.Println("GetProperty pID error ", err)
		request.GetConnection().Stop()
		return
	}

	fmt.Printf("player pID = %d , move(%f,%f,%f,%f)\n", pID, msg.X, msg.Y, msg.Z, msg.V)
	// 3. 根据pID得到player对象
	player := core.WorldMgrObj.GetPlayerByPID(pID.(int32))

	// 4. 让player对象发起移动位置信息广播
	player.UpdatePos(msg.X, msg.Y, msg.Z, msg.V)
}
