package znet

import (
	"fmt"
	"go-study/zinx/ziface"
	"strconv"
)

// MsgHandler 消息处理模块的实现
type MsgHandle struct {
	//  存放每个MsgId对应的处理方法
	Apis map[uint32]ziface.IRouter
}

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis: make(map[uint32]ziface.IRouter),
	}
}
func (mh *MsgHandle) DoMsgHandler(request ziface.IRequest) {
	// 1 从request中找到msgId
	handler, ok := mh.Apis[request.GetMsgId()]
	if !ok {
		fmt.Println("api msgId = ", request.GetMsgId(), " is NOT FOUND! Need Register!")
	}
	// 2 根据MsgId 调度router对应的业务
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

// 为消息添加具体的处理逻辑
func (mh *MsgHandle) AddRouter(msgId uint32, router ziface.IRouter) {
	// 1 判断 当前msg绑定的API处理方法是否已经存在
	if _, ok := mh.Apis[msgId]; ok {
		panic("repeat api, msg Id = " + strconv.Itoa(int(msgId)))
	}
	// 2 添加msg和API的绑定关系
	mh.Apis[msgId] = router
	fmt.Println("Add api MsgId = ", msgId, "succ!")
}
