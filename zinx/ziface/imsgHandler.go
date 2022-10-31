package ziface

// IMsgHandler 消息管理抽象层
type IMsgHandle interface {
	// 调度/执行对应的Router消息处理方法
	DoMsgHandler(request IRequest)

	// 为消息添加具体的处理逻辑
	AddRouter(msgId uint32, router IRouter)

	// 启动 Worker 工作池
	StartWorkerPool()

	// 将消息交给 TaskQueue, 由worker 进行处理
	SendMsgToTaskQueue(request IRequest)
}
