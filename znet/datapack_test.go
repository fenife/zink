package znet

import (
	"fmt"
	"io"
	"net"
	"testing"
	"time"
)

const testAddr = "127.0.0.1:7777"

//只是负责测试datapack拆包 封包的单元测试
func TestDataPack(t *testing.T) {
	//模拟的服务器
	//1 创建socketTCP
	listener, err := net.Listen("tcp", testAddr)
	if err != nil {
		fmt.Println("serve listen err:", err)
	}

	//创建一个go承载负责从客户端处理业务
	go func() {
		//2.c从客户端读取数据，拆包处理
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("server accept error", err)
			}
			go func(conn net.Conn) {
				//处理客户端请求
				//拆包的过程
				//定义一个拆包的对象dp
				dp := NewDataPack()
				for {
					//1.第一次从conn读 把包的head读出来
					headData := make([]byte, dp.GetHeadLen())
					if _, err := io.ReadFull(conn, headData); err != nil {
						fmt.Println("read head err:", err)
						break
					}
					msgHead, err := dp.Unpack(headData)
					if err != nil {
						fmt.Println("server unpack err:", err)
						break
					}
					if msgHead.GetMsgLen() > 0 {
						//msg是有数据的，需要进行第二次读取
						//2.第二次读从conn读，根据head中dataLen再读取data内容
						msg := msgHead.(*Message)
						msg.Data = make([]byte, msg.GetMsgLen())
						_, err := io.ReadFull(conn, msg.Data)
						if err != nil {
							fmt.Println("server unpack data err:", err)
							break
						}
						//完整的一个消息已经读取完毕
						fmt.Printf("---> Recv MsgId: %d, dataLen = %d, data = %s\n",
							msg.Id, msg.DataLen, msg.Data)
					}
				}
			}(conn)
		}
	}()

	//模拟客户端
	conn, err := net.Dial("tcp", testAddr)
	if err != nil {
		fmt.Println("client dial err:", err)
		return
	}
	//创建一个封包对象 dp
	dp := NewDataPack()

	//模拟粘包过程，封装两个msg一同发送
	//封装第一个msg1包
	msg1 := &Message{
		Id: 1,
		DataLen: 4,
		Data: []byte("zinx"),
	}
	sendData1, err := dp.Pack(msg1)
	if err != nil {
		fmt.Println("client pack msg1 error", err)
		return
	}
	//封装第二个msg2包
	msg2 := &Message{
		Id: 2,
		DataLen: 7,
		Data: []byte("hello!!"),
	}
	sendData2, err := dp.Pack(msg2)
	if err != nil {
		fmt.Println("client pack msg2 error", err)
		return
	}

	//将2个包粘在一起
	sendData1 = append(sendData1, sendData2...)
	for {
		//一次性发送给服务端
		conn.Write(sendData1)
		time.Sleep(3 * time.Second)
	}
}

// todo: 不用tcp连接，直接构建数据包测试？