package ziface

/**
 消息管理抽象层
 */

type IMsgHandle interface {
	DoMsgHandler(request IRequest)          //马上以非阻塞方式处理消息
	AddRouter(msgId uint32, router IRouter) //注册路由
}
