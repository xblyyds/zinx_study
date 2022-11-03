package znet

import (
	"errors"
	"fmt"
	"go-study/zinx/ziface"
	"sync"
)

/**
连接管理模块
*/
type ConnManager struct {
	// 管理的连接集合
	connections map[uint32]ziface.IConnection
	// 保护连接的读写锁
	connLock sync.RWMutex
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]ziface.IConnection),
	}
}

func (connMgr *ConnManager) Add(conn ziface.IConnection) {
	// 保护共享资源 map,加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	// 添加进集合
	connMgr.connections[conn.GetConnID()] = conn

	fmt.Println("connID = ", conn.GetConnID(), " add to ConnManager successfully: conn num = ", connMgr.Len())
}

func (connMgr *ConnManager) Remove(conn ziface.IConnection) {
	// 保护共享资源 map,加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	delete(connMgr.connections, conn.GetConnID())

	fmt.Println("connID = ", conn.GetConnID(), " remove from ConnManager successfully: conn num  = ", connMgr.Len())
}

func (connMgr *ConnManager) Get(connID uint32) (ziface.IConnection, error) {
	// 保护共享资源 map,加写锁
	connMgr.connLock.RLock()
	defer connMgr.connLock.RUnlock()

	if conn, ok := connMgr.connections[connID]; ok {
		return conn, nil
	} else {
		return nil, errors.New("connection not FOUND")
	}
}

func (connMgr *ConnManager) Len() int {
	return len(connMgr.connections)
}

func (connMgr *ConnManager) ClearConn() {
	// 保护共享资源 map,加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	// 删除 并停止 conn
	for connID, conn := range connMgr.connections {
		// 停止
		conn.Stop()
		// 删除
		delete(connMgr.connections, connID)
	}

	fmt.Println("Clear All connections succ! conn num = ", connMgr.Len())
}