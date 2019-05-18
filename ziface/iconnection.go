package ziface

import (
	"net"
)

type IConnection interface {
	//启动连接，让当前连接工作
	Start()
	//停止连接，结束当前连接状态
	Stop()
	//从当前连接获取原始的socket TCPConn
	GetTCPConnection() *net.TCPConn
	//获取当前连接ID
	GetConnID() uint32
	//获取远程客户端地址信息
	RemoteAddr() net.Addr
	//直接将Message数据发送给远程的TCP客户端（无缓冲）
	SendMsg(msgId uint32, data []byte) error
	//直接将Message数据发送给远程的TCP客户端（有缓冲）
	SendBuffMsg(msgId uint32, data []byte) error

	//设置连接属性
	SetProperty(key string, value interface{})
	//获取连接属性
	GetProperty(key string) (interface{}, error)
	//移除连接属性
	RemoveProperty(key string)
}

//定义一个统一处理连接的业务接口 参数（原生TCP连接、客户端请求的数据、客户端请求的数据长度）
type HandFunc func(*net.TCPConn, []byte, int) error
