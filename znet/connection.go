package znet

import (
	"net"
	"zink/ziface"
)

/*
 链接模块
*/
type Connection struct {
	//当前链接的socket TCP套接字
	Conn *net.TCPConn
	//链接的ID
	ConnID uint32
	//当前的链接状态
	isClosed bool
	//当前链接所绑定的业务
	handleAPI ziface.HandleFunc
	//告知当前链接已经退出的/停止 channel
	ExitChan chan bool
}

//初始化链接模块的方法
func NewConnection(conn *net.TCPConn, connID uint32, callback ziface.HandleFunc) *Connection {
	c := &Connection{
		Conn:      conn,
		ConnID:    connID,
		handleAPI: callback,
		isClosed:  false,
		ExitChan:  make(chan bool, 1),
	}

	return c
}
