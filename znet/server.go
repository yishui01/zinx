package znet

import (
	"fmt"
	"github.com/pkg/errors"
	"net"
	"time"
	"zinx/utils"
	"zinx/ziface"
)

/**
流程总结
服务端开启tcp服务端，listen+循环accept
accept进来一个客户端，生成一个客户端结构体，生成结构体的时候会把处理这个客户端的handler传进去
再调用客户端结构体的start方法，start再调用statReader，顺便监听结构体中的结束channel是否有消息，
startReader中再调用handler处理函数，出现异常或者客户端数据处理完毕时向结束channel发送消息
 */

/**
V2.0
将客户端处理函数变为路由功能，server结构体添加路由，然后创建客户端struct的时候，将路由传到里面，有客户端协程执行handler
 */

//iServer 接口实现，定义一个Server服务端
type Server struct {
	//服务器的名称
	Name string
	//tcp4 or other
	IPVersion string
	//服务绑定的IP地址
	IP string
	//服务器绑定的端口
	Port int
	//当前Server由用户绑定的回调router，也就是Server注册的连接对应的处理业务
	Router ziface.IRouter
}

//============================ 定义当前客户端连接的handle api ==============================//
func CallBackToClient(conn *net.TCPConn, data []byte, cnt int) error {
	//回显业务
	fmt.Println("[连接handler处理函数] CallBackToClient ... ")
	if _, err := conn.Write(data[:cnt]); err != nil {
		fmt.Println("write back buf err ", err)
		return errors.New("处理函数CallBackToClient发生故障")
	}
	return nil
}

//============================= 实现 ziface.IServer 里的全部接口方法=======================//

//开启网络服务器
func (s *Server) Start() {
	fmt.Printf("[START] Server listener at IP: %s, Port %d, is starting\n", s.IP, s.Port)
	fmt.Printf("[Zinx] Version: %s, MaxConn: %d, MaxPacketSize: %d\n", utils.G_Obj.Version, utils.G_Obj.MaxConn, utils.G_Obj.MaxPacketSize)

	//开启一个go去做服务端Listener业务
	go func() {
		//1、获取一个TCP的Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr err: ", err)
			return
		}

		//2、监听服务器地址
		listenner, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen", s.IPVersion, "ERR:", err)
			return
		}

		//已经监听成功
		fmt.Println("start Zinx server ", s.Name, " succ, now listenning...")

		//TODO server.go 应该有一个自动生成ID的方法
		var cid uint32
		cid = 0

		//3、启动server网络连接业务
		for {
			//3.1 阻塞等待客户端建立连接请求
			conn, err := listenner.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err", err)
				continue
			}

			//3.2 TODO Server.Start() 设置服务器最大连接控制，如果超过最大连接，那么则关闭此新的连接

			//3.3 TODO Server.Start() 处理该新连接请求的 业务 方法， 此时应该有handler 和 conn是绑定的
			dealConn := NewConnection(conn, cid, s.Router)
			cid++

			//3.4 启动当前连接的处理业务
			go dealConn.Start()

		}

	}()

}

func (s *Server) Stop() {
	fmt.Println("[STOP] Zinx server , name ", s.Name)

	//TODO Server.Stop()  將其他需要清理的连接信息或者其他信息　也要一并停止或者清理
}

func (s *Server) Serve() {
	s.Start()

	//TODO Server.Serve() 是否在启动服务的时候 还要处理其他的事情呢 可以在这里添加

	//阻塞，否则主Go退出， listener的go将会退出
	for {
		time.Sleep(10 * time.Second)
	}
}

//路由功能：给当前服务注册一个路由业务方法，供客户端连接处理使用
func (s *Server) AddRouter(router ziface.IRouter) {
	s.Router = router
	fmt.Println("Add Router success! ")
}

func NewServer(name string) ziface.IServer {
	//先初始化全局配置文件
	err := utils.G_Obj.Reload()
	var s *Server
	if err != nil {
		fmt.Printf("采用默认Zinx配置")
		s = &Server{
			Name:      "默认Zinx Server",
			IPVersion: "tcp4",
			IP:        "127.0.0.1",
			Port:      7777,
			Router:    nil,
		}
	} else {
		s = &Server{
			Name:      utils.G_Obj.Name,
			IPVersion: "tcp4",
			IP:        utils.G_Obj.Host,
			Port:      utils.G_Obj.TcpPort,
			Router:    nil,
		}
	}

	return s
}
