package main

import "zink/znet"

/*
 基于Zinx框架来开发的 服务器端应用程序
*/
func main() {
	// 1.创建一个server句柄，使用Zinx的api
	s := znet.NewServer("[zink v0.2]")
	// 2.启动server
	s.Serve()
}
