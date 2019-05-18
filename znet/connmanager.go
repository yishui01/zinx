package znet

import (
	"fmt"
	"github.com/pkg/errors"
	"sync"
	"zinx/ziface"
)

//连接管理结构体
type ConnManager struct {
	connections map[uint32]ziface.IConnection
	connLock    sync.RWMutex
}

/**
	创建一个连接管理结构体
 */
func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]ziface.IConnection),
	}
}

/**
	添加连接
 */
func (connMgr *ConnManager) Add(conn ziface.IConnection) {
	//保护共享资源Map 加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	//将conn连接添加到ConnManager中
	connMgr.connections[conn.GetConnID()] = conn

	fmt.Println("connection add to ConnManager successfully: conn num = ", connMgr.Len())
}

//删除连接
func (connMgr *ConnManager) Remove(conn ziface.IConnection) {
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	//删除连接信息
	delete(connMgr.connections, conn.GetConnID())

	fmt.Println("connection Remove ConnID=", conn.GetConnID(), "successfully: conn num = ", connMgr.Len())
}

//获取连接
func (connMgr *ConnManager) Get(connID uint32) (ziface.IConnection, error) {
	//加锁
	connMgr.connLock.RLock()
	defer connMgr.connLock.Unlock()

	if conn, ok := connMgr.connections[connID]; ok {
		return conn, nil
	} else {
		return nil, errors.New("connection not found")
	}
}

func (connMgr *ConnManager) Len() int {
	return len(connMgr.connections)
}

func (connMgr *ConnManager) ClearConn() {
	//保护共享资源Map，加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	//停止并删除全部的连接消息
	for connID, conn := range connMgr.connections {
		//停止
		conn.Stop()
		delete(connMgr.connections, connID)
	}

	fmt.Println("Clear All Connections successfully: conn nun = ", connMgr.Len())
}
