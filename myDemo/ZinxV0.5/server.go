package main

import (
	"fmt"
	"zink/ziface"
	"zink/znet"
)

/*
 基于Zinx框架来开发的 服务器端应用程序
*/

//ping test 自定义路由
type PingRouter struct {
	znet.BaseRouter
}


// Test Handle
func (br *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Router Handle...")
	//读取客户端的数据，再回写ping
	fmt.Printf("recv from client: msgId = %d, data = %s",
		request.GetMsgID(), request.GetData())
	err := request.GetConnection().SendMsg(1, []byte("ping...ping...ping\n"))
	if err != nil {
		fmt.Println(err)
	}
}


func main() {
	// 1.创建一个server句柄，使用Zinx的api
	s := znet.NewServer("[zink v0.3]")
	//2.给当前zinx添加一个自定义的router
	s.AddRouter(&PingRouter{})
	// 3.启动server
	s.Serve()
}
