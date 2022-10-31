package znet

import (
	"fmt"
	"go-study/zinx/utils"
	"go-study/zinx/ziface"
	"strconv"
)

// MsgHandler 消息处理模块的实现
type MsgHandle struct {
	//  存放每个MsgId对应的处理方法
	Apis map[uint32]ziface.IRouter
	// 负责 Worker 取任务的消息队列
	TaskQueue []chan ziface.IRequest
	// 业务工作 Worker 池的worker 数量
	WorkPoolSize uint32
}

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis:         make(map[uint32]ziface.IRouter),
		WorkPoolSize: utils.GlobalObject.WorkerPoolSize,
		TaskQueue:    make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize),
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

// 启动一个 Worker 工作池
func (mh *MsgHandle) StartWorkerPool() {
	// 根据 workerPoolSize 分别开启 worker, 每个 worker 用一个 goroutine 承载
	for i := 0; i < int(mh.WorkPoolSize); i++ {
		// 一个 worker启动
		// 1 当前的 worker对应的 channel消息队列, 开辟空间 第0个 worker 就用 第 0 个channel...
		mh.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		// 启动当前的worker， 阻塞等待消息从 channel 传递进来
		go mh.StartOneWorker(i, mh.TaskQueue[i])
	}
}

// 启动一个 Worker 工作流程
func (mh *MsgHandle) StartOneWorker(workID int, taskQueue chan ziface.IRequest) {
	fmt.Println("WorkerID = ", workID, " is started...")
	// 不断的阻塞等待对应的消息队列的消息
	for {
		select {
		case request := <-taskQueue:
			mh.DoMsgHandler(request)
		}
	}
}

// 将消息交给TaskQueue, 由worker进行处理
func (mh *MsgHandle) SendMsgToTaskQueue(request ziface.IRequest) {
	// 1 将消息分配给不同的worker
	// 根据客户端建立的ConnID来分配
	workerId := request.GetConnection().GetConnID() % mh.WorkPoolSize
	fmt.Println("Add ConnID = ", request.GetConnection().GetConnID(), ", request MsgId = ", request.GetMsgId(), " to WorkerID = ", workerId)
	// 2 将消息发送给对应的 worker 的 TaskQueue
	mh.TaskQueue[workerId] <- request
}
