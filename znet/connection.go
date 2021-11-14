package znet

import (
	"fmt"
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

//链接的读业务方法
func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running...")
	defer fmt.Printf("ConnID = %d, Reader is exit, remote addr is %s\n",
		c.ConnID, c.RemoteAddr().String())
	defer c.Stop()

	for {
		//读客户端的数据到buf中，最大字节512
		buf := make([]byte, 512)
		cnt, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("recv buf err", err)
			break
		}
		//调用当前链接所绑定的HandleAPI
		if err := c.handleAPI(c.Conn, buf, cnt); err != nil {
			fmt.Printf("ConnID %d handle is error: %s\n", c.ConnID, err)
			break
		}
	}
}


//启动链接 让当前的链接准备开始工作
func (c *Connection) Start() {
	fmt.Println("Conn Start() ... ConnID =", c.ConnID)

	//启动从当前链接的读数据的业务
	go c.StartReader()

	//TODO 启动从当前链接写数据的业务
}

//停止链接 结束当前链接的工作
func (c *Connection) Stop() {
	fmt.Println("Conn Stop().. ConnID =", c.ConnID)
	//如果当前链接已经关闭
	if c.isClosed {
		return
	}
	//关闭socket链接
	c.Conn.Close()
	//回收资源
	close(c.ExitChan)
}

//获取当前链接的绑定socket conn
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

//获取当前链接模块的链接ID
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

//获取远程客户端的 TCP状态 IP port
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

//发送数据，将数据发送给远程的客户端
func (c *Connection) Send(data []byte) error {
	return nil
}
