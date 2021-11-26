package znet

import (
	"fmt"
	"zink/ziface"
)

/*
消息处理模块的实现
*/

type MsgHandler struct {
	//存放每个MsgID 所对应的处理方法
	Apis map[uint32]ziface.IRouter
}

//初始化创建MsgHandler方法
func NewMsgHandler() *MsgHandler {
	return &MsgHandler{
		Apis: make(map[uint32]ziface.IRouter),
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

