package znet

import (
	"fmt"
	"net"
	"testing"
	"time"
)

func ClientTest() {
	fmt.Println("Client Test ... start")
	//3秒之后发起测试请求， 给服务端开启服务端机会
	time.Sleep( 5 * time.Second)

	conn,err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("client start err, exit,err is : ",err)
		return
	}

	for {
		_, err := conn.Write([]byte("hello Zinx"))
		if err != nil {
			fmt.Println("Write error err", err)
			return
		}

		buf := make([]byte, 512)
		cnt, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read buf err", err)
			return
		}

		fmt.Printf("server call back : %s, cnt = %d\n", buf, cnt)

		time.Sleep(1 * time.Second)
	}
}


//Server 模块的测试函数
func TestServer(t *testing.T) {

	s := NewServer("[zinx V0.1]")

	//进行客户端测试
	go ClientTest()

	//2 开启服务
	s.Serve()
}	
