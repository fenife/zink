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

type HelloZinxRouter struct {
	znet.BaseRouter
}


// Test Handle
func (br *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call PingRouter Handle...")
	//读取客户端的数据，再回写ping
	fmt.Printf("recv from client: msgId = %d, data = %s\n",
		request.GetMsgID(), request.GetData())
	err := request.GetConnection().SendMsg(request.GetMsgID(), []byte("ping router"))
	if err != nil {
		fmt.Println(err)
	}
}

func (br *HelloZinxRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call HelloZinxRouter Handle...")
	//读取客户端的数据，再回写ping
	fmt.Printf("recv from client: msgId = %d, data = %s\n",
		request.GetMsgID(), request.GetData())
	err := request.GetConnection().SendMsg(request.GetMsgID(), []byte("hello zinx router"))
	if err != nil {
		fmt.Println(err)
	}
}


func main() {
	// 1.创建一个server句柄，使用Zinx的api
	s := znet.NewServer("[zink v0.6]")
	//2.给当前zinx添加自定义的router
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloZinxRouter{})
	// 3.启动server
	s.Serve()
}
