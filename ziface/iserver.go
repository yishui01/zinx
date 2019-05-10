package ziface

//定义服务器接口
type Iserver interface {
	//启动服务器方法
	Start()
	//停止服务器方法
	Stop()
	//开启业务方法
	Serve()
}
