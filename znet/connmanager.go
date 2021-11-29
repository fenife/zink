package znet

import (
	"fmt"
	"sync"
	"zink/ziface"
)

//链接管理模块
type ConnManager struct {
	//管理的链接集合
	connections map[uint32]ziface.IConnection
	//保护链接集合的读写锁
	connLock sync.RWMutex
}

//创建当前链接的方法
func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]ziface.IConnection),
	}
}

//添加链接
func (m *ConnManager) Add(conn ziface.IConnection) {
	m.connLock.Lock()
	defer m.connLock.Unlock()

	m.connections[conn.GetConnID()] = conn
	fmt.Printf("connID = %d add to ConnManager successfully, conn num = %d",
		conn.GetConnID(), m.Len())
}

//删除链接
func (m *ConnManager) Remove(conn ziface.IConnection) {
	m.connLock.Lock()
	defer m.connLock.Unlock()

	delete(m.connections, conn.GetConnID())
	fmt.Printf("connID = %d remove from ConnManager, conn num = %d",
		conn.GetConnID(), m.Len())
}

//根据connID获取链接
func (m *ConnManager) Get(connID uint32) (ziface.IConnection, error) {
	m.connLock.Lock()
	defer m.connLock.Unlock()

	if conn, ok := m.connections[connID]; ok {
		return conn, nil
	} else {
		return nil, fmt.Errorf("connID = %d not FOUND", connID)
	}
}

//得到当前链接总数
func (m *ConnManager) Len() int {
	return len(m.connections)
}

//清楚并终止所有链接
func (m *ConnManager) ClearConn() {
	m.connLock.Lock()
	defer m.connLock.Unlock()

	for connID, conn := range m.connections {
		conn.Stop()

		delete(m.connections, connID)
	}
	fmt.Printf("clear all connnections succ! conn num = %d", m.Len())
}

