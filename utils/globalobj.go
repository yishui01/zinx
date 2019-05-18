package utils

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"zinx/ziface"
)

/**
 全局配置参数
 */

type GlobalObj struct {
	/**
		 Server
	 */
	TcpServer ziface.IServer //当前Zinx的全局Server对象
	Host      string         //当前服务器主机IP
	TcpPort   int            //当前服务器主机监听端口号
	Name      string         //当前服务器名称

	/**
		 Zinx
	 */
	Version          string //当前Zinx版本号
	MaxPacketSize    uint32 //读取数据包的最大值
	MaxConn          int    //当前服务器主机允许的最大连接个数
	WorkerPoolSize   uint32 //业务工作Worker池数量
	MaxWorkerTaskLen uint32 //业务工作Worker对应负责的任务队列最大任务存储数量
	MaxMsgChanLen    uint32
	/**
		config file path
	 */
	ConfigFilePath string
}

var G_Obj *GlobalObj

//读取用户配置文件
func (g *GlobalObj) Reload() error {
	configFileName := "conf/zinx.json"
	_, err := os.Stat(configFileName)
	if err != nil {
		return err
	}
	data, err := ioutil.ReadFile(configFileName)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &G_Obj)
	if err != nil {
		return err
	}
	return nil
}

func init() {
	G_Obj = &GlobalObj{
		Name:             "Zinx App Server",
		Version:          "V0.4",
		TcpPort:          7777,
		Host:             "0.0.0.0",
		MaxConn:          12000,
		MaxPacketSize:    4096,
		ConfigFilePath:   "conf/zinx.json",
		WorkerPoolSize:   10,
		MaxWorkerTaskLen: 1024,
		MaxMsgChanLen:    1024, //goroutine之间通过有缓冲管道通信时管道的最大长度
	}

	//加载用户自定义的配置文件，覆盖默认配置
	G_Obj.Reload()
}
