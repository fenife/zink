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

//创建链接之后执行钩子函数
func DoConnectionBegin(conn ziface.IConnection) {
	fmt.Println("===> DoConnectionBegin is called")
	err := conn.SendMsg(202, []byte("DoConnection Begin"))
	if err != nil {
		fmt.Println(err)
	}

	//给当前的链接设置一些属性
	fmt.Println("set conn property ...")
	conn.SetProperty("Name", "Zink")
}

//链接端口之前的需要执行的函数
func DoConnectionLost(conn ziface.IConnection) {
	fmt.Println("===> DoConnectionLost is called")
	fmt.Printf("conn ID = %d is lost\n", conn.GetConnID())

	//获取链接属性
	if name, err := conn.GetProperty("Name"); err == nil {
		fmt.Println("Name = ", name)
	}
}

func main() {
	// 1.创建一个server句柄，使用Zinx的api
	s := znet.NewServer("[zink v1.0]")

	//注册链接hook钩子函数
	s.SetOnConnStart(DoConnectionBegin)
	s.SetOnConnStop(DoConnectionLost)

	//2.给当前zinx添加自定义的router
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloZinxRouter{})
	// 3.启动server
	s.Serve()
}
