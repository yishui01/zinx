package znet

import (
	"fmt"
	"github.com/pkg/errors"
	"io"
	"net"
	"zinx/utils"
	"zinx/ziface"
)

type Connection struct {
	//当前连接的socket TCP套接字
	Conn *net.TCPConn
	//当前连接的ID 也可以称作为SessionID、ID全局唯一
	ConnID uint32
	//当前连接的关闭状态
	isClosed bool

	//该连接的处理方法API
	Msghandler ziface.IMsgHandle

	//告知该连接已经退出/停止的channel
	ExitBuffChan chan bool

	//无缓冲管道，用于读、写两个goroutine之间的消息通道
	msgChan chan []byte
}

//创建连接的方法
func NewConnection(conn *net.TCPConn, connID uint32, handler ziface.IMsgHandle) *Connection {
	c := &Connection{
		Conn:         conn,
		ConnID:       connID,
		isClosed:     false,
		Msghandler:   handler,
		ExitBuffChan: make(chan bool, 1),
		msgChan:      make(chan []byte),
	}

	return c
}

//启动连接，让当前连接开始工作
func (c *Connection) Start() {
	//开启处理该链接读取到客户端数据之后的请求业务
	go c.StartReader()
	go c.StartWriter()

	for {
		select {
		case <-c.ExitBuffChan:
			//得到退出消息,不再阻塞
			return
		}
	}
}

//停止连接，结束当前连接状态
func (c *Connection) Stop() {
	//1、如果当前连接已关闭
	if c.isClosed == true {
		return
	}
	c.isClosed = true

	//TODO Connection Stop() 如果用户注册了该链接的关闭回调业务，那么此刻应该显示调用

	// 关闭socket连接
	c.Conn.Close()

	//通知从缓冲队列读数据的业务，该链接已经关闭
	c.ExitBuffChan <- true

	//关闭该链接全部管道
	close(c.ExitBuffChan)
}

//从当前连获取原始的socket	 TCPConn
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

//获取当前连接ID
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

//获取远程客户端地址信息
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

//处理conn读取数据的Goroutine
func (c *Connection) StartReader() {

	fmt.Println("Reader Goroutine is running")
	defer fmt.Println(c.RemoteAddr().String(), " conn reader exit!")
	defer c.Stop()

	for {
		//创建解包对象
		dp := NewDataPack()

		//读取客户端的Msg head
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("read msg head error ", err)
			c.ExitBuffChan <- true
			continue

		}
		//拆包，得到msgid 和 datalen 放在msg中
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack head error", err)
			c.ExitBuffChan <- true
			continue
		}

		//根据 dataLen 读取 data，放在msg.Data中
		var data []byte
		if msg.GetDataLen() > 0 {
			data = make([]byte, msg.GetDataLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("read msg data error ", err)
				c.ExitBuffChan <- true
				continue
			}
		}
		msg.SetData(data)
		//得到当前客户端请求的Request数据
		req := Request{
			conn: c,
			msg:  msg,
		}
		if utils.G_Obj.WorkerPoolSize > 0 {
			//已经启动工作池机制，将消息交给Worker处理
			c.Msghandler.SendMsgToTaskQueue(&req)
		} else {
			//从绑定好的消息和对应的处理方法中执行对应的Handle方法
			go c.Msghandler.DoMsgHandler(&req)
		}

		//从路由Routers 中找到注册绑定的Conn对应的Handle
		//		go c.Msghandler.DoMsgHandler(&req)

	}
}

//将消息封包
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("Connection closed when send msg")
	}
	dp := NewDataPack()
	msg, err := dp.Pack(NewMsg(msgId, data))
	if err != nil {
		fmt.Println("Pack error msg id = ", msgId)
		return err
	}
	//发送给客户端
	c.msgChan <- msg

	return nil
}

func (c *Connection) StartWriter() {
	defer fmt.Println(c.RemoteAddr().String(), " conn Writer exit!")

	for {
		select {
		case data := <-c.msgChan:
			//有数据要写给客户端
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send Data error:", err, "Conn Writer exit")
				return
			}
		case <-c.ExitBuffChan:
			//Conn已关闭
			return
		}
	}
}
