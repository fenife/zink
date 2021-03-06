package znet

import (
	"fmt"
	"zink/utils"
	"zink/ziface"
)

/*
消息处理模块的实现
*/

type MsgHandler struct {
	//存放每个MsgID 所对应的处理方法
	Apis map[uint32]ziface.IRouter
	//负责Worker取任务的消息队列
	TaskQueue []chan ziface.IRequest
	//负责工作Worker池的worker数量
	WorkerPoolSize uint32
}

//初始化创建MsgHandler方法
func NewMsgHandler() *MsgHandler {
	return &MsgHandler{
		Apis:           make(map[uint32]ziface.IRouter),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize, //从全局配置中获取
		TaskQueue:      make([]chan ziface.IRequest, int(utils.GlobalObject.WorkerPoolSize)),
	}
}

//执行对应的Router消息处理方法
func (h *MsgHandler) DoMsgHandler(request ziface.IRequest) {
	//1 从 Request中找到msgID
	handler, ok := h.Apis[request.GetMsgID()]
	if !ok {
		panic(fmt.Sprintf("api msgID = %v is not found! need registered",
			request.GetMsgID()))
	}
	//2 根据MsgID 调度对应的router业务极客
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

//为消息添加具体的处理逻辑
func (h *MsgHandler) AddRouter(msgID uint32, router ziface.IRouter) {
	// 判断 当前msg绑定的API处理方法是否已经存在
	if _, ok := h.Apis[msgID]; ok {
		// id已经注册了
		panic(fmt.Sprintf("repead api, msgID = %d", msgID))
	}
	h.Apis[msgID] = router
	fmt.Printf("Add api msgID = %d succ\n", msgID)
}

//启动一个Worker工作池(开启工作池的动作只能发生一次，一个zinx框架只能有一个worker工作池)
func (h *MsgHandler) StartWorkerPool() {
	//根据WorkerPoolSize分别开启Worker，每个Worker用一个go来承载
	for i := 0; i < int(h.WorkerPoolSize); i++ {
		//一个Worker被启动
		//1 当前的worker对应的channel消息队列 开辟空间
		h.TaskQueue[i] = make(chan ziface.IRequest, int(utils.GlobalObject.MaxWorkerTaskLen))
		//2 启动当前的Worker，阻塞等待消息从channel传递过来
		go h.StartOneWorker(i, h.TaskQueue[i])
	}
}

func (h *MsgHandler) StartOneWorker(workerID int, taskQueue chan ziface.IRequest) {
	fmt.Printf("worker id = %d is started ...\n", workerID)
	//不断的阻塞等待对应消息队列的消息
	for {
		select {
		//如果有消息过来，出列的就是一个客户端的Request, 执行当前Request所绑定的业务
		case request := <-taskQueue:
			h.DoMsgHandler(request)
		}
	}
}

//将消息交给TaskQueue，由worker处理
func (h *MsgHandler) SendMsgToTaskQueue(request ziface.IRequest) {
	//将消息评价分配给不同的worker
	//根据客户端简历的ConnID来进行分配
	workerID := request.GetConnection().GetConnID() % h.WorkerPoolSize

	fmt.Printf("Add ConnID=%d request msgID=%d to workerID=%d\n",
		request.GetConnection().GetConnID(), request.GetMsgID(), workerID)

	//将消息发送给对应的worker的TaskQueue即可
	h.TaskQueue[workerID] <- request
}

