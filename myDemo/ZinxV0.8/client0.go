package main

import (
	"fmt"
	"io"
	"net"
	"time"
	"zink/znet"
)

/*
 模拟客户端
*/
func main() {
	fmt.Println("client0 start ...")
	time.Sleep(1 * time.Second)
	// 1.直接连接远程服务器，得到一个conn连接
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("client start err, exit!")
		return
	}

	//2 连接调用Write 写数据
	for {
		//发送封包的message消息
		dp := znet.NewDataPack()
		binaryMsg, err := dp.Pack(znet.NewMsgPackage(0, []byte("zinx client0 Test Message")))
		if err != nil {
			fmt.Println("pack error", err)
			return
		}

		_, err = conn.Write(binaryMsg)
		if err != nil {
			fmt.Println("write error", err)
			return
		}

		//服务器应该给我们回复一个message数据， MsgID: 0, pingpingping
		//先读取流中的head部分 得到ID 和 dataLen
		binaryHead := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(conn, binaryHead); err != nil {
			fmt.Println("read head error", err)
			break
		}

		//将二进制的head拆包到msg结构体中
		msgHead, err := dp.Unpack(binaryHead)
		if err != nil {
			fmt.Println("client unpack msgHead error", err)
			break
		}
		if msgHead.GetMsgLen() > 0 {
			//再根据DataLen进行第二次读取，将data读出来
			msg := msgHead.(*znet.Message)
			msg.Data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(conn, msg.Data); err != nil {
				fmt.Println("read msg data error", err)
				return
			}
			fmt.Printf("---> recv server msg: id = %d, len = %d, data = %s\n",
				msg.Id, msg.DataLen, msg.Data)

		}

		// cpu 阻塞
		time.Sleep(1 * time.Second)
	}

}
