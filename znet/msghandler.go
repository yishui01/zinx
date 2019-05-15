package znet

import (
	"fmt"
	"strconv"
	"zinx/ziface"
)

type MsgHandle struct {
	Apis map[uint32]ziface.IRouter //全局路由map
}

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis: make(map[uint32]ziface.IRouter),
	}
}

//执行消息对应的路由处理方法
func (mh *MsgHandle) DoMsgHandler(request ziface.IRequest) {
	handler, ok := mh.Apis[request.GetMsgID()]
	if !ok {
		fmt.Println("api msgID = ", request.GetMsgID(), "is not Found in map")
		return
	}

	//执行对应的方法
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

//添加路由
func (mh *MsgHandle) AddRouter(msgId uint32, router ziface.IRouter) {
	//1、判断当前，msgID是否已经绑定了处理方法
	if _, ok := mh.Apis[msgId]; ok {
		panic("repeat api , msgID = " + strconv.Itoa(int(msgId)))
	}
	//2、添加msg与api的关系
	mh.Apis[msgId] = router
	fmt.Println("Add api msgId = ", msgId)
}
