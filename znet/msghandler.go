package znet

import (
	"fmt"
	"strconv"
	"zinx/utils"
	"zinx/ziface"
)

type MsgHandle struct {
	Apis           map[uint32]ziface.IRouter //全局路由map
	WorkerPoolSize uint32                    //业务工作worker池的数量
	TaskQueue      []chan ziface.IRequest    //worker负责取任务的消息队列
}

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis:           make(map[uint32]ziface.IRouter),
		WorkerPoolSize: utils.G_Obj.WorkerPoolSize,
		TaskQueue:      make([]chan ziface.IRequest, utils.G_Obj.WorkerPoolSize),
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

//启动一个Worker工作流程
func (mh *MsgHandle) StartOneWorker(workerID int, taskQueue chan ziface.IRequest) {
	fmt.Println("Worker ID = ", workerID, "is started")
	//不断的等待队列中的消息
	for {
		select {
		//有消息则取出执行
		case request := <-taskQueue:
			mh.DoMsgHandler(request)
		}
	}
}

//启动worker工作池
func (mh *MsgHandle) StartWorkerPool() {
	//遍历需要启动worker数量，依次启动
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		//一个worker被启动
		//初始化队列
		mh.TaskQueue[i] = make(chan ziface.IRequest, utils.G_Obj.MaxWorkerTaskLen)
		//启动当前Worker，阻塞的等待对应的任务队列是否有消息传递进来
		go mh.StartOneWorker(i, mh.TaskQueue[i])
	}
}

//将消息交给TaskQueue，由worker进行处理
func (mh *MsgHandle) SendMsgToTaskQueue(request ziface.IRequest) {
	//根据ConnId取模来分配当前连接到不同的队列中,要保证同一个用户的连接要被分配到同一个队列中
	workerID := request.GetConnection().GetConnID() % mh.WorkerPoolSize
	fmt.Println("Add ConnID=", request.GetConnection().GetConnID(), "request msg id=", request.GetMsgID(), "to worker --", workerID)
	mh.TaskQueue[workerID] <- request
}
