package main

import (
	"fmt"
	"go-study/zinx/znet"
	"io"
	"net"
	"time"
)

func main() {
	fmt.Println("client0 start...")

	time.Sleep(1 * time.Second)

	// 1、直接连接远程服务器, 得到一个conn连接
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("client start err, exit!")
		return
	}

	for {
		// 发送封包的 Msg
		dp := znet.NewDataPack()
		binaryMsg, err := dp.Pack(znet.NewMessage(0, []byte("ZinxV0.6 client0 Test Message")))
		if err != nil {
			fmt.Println("Pack error ", err)
			return
		}
		if _, err := conn.Write(binaryMsg); err != nil {
			fmt.Println("write error: ", err)
			return
		}
		// 服务器应该回复一个 message数据 msgId: 1 data: ping...ping...ping

		// 1 先读取head部分, 得到msgId 和 dataLen
		binaryHead := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(conn, binaryHead); err != nil {
			fmt.Println("read head error: ", err)
			break
		}

		// 将二进制的head拆包到msg结构体中
		msgHead, err := dp.UnPack(binaryHead)
		if err != nil {
			fmt.Println("client unpack msgHead error: ", err)
			break
		}
		if msgHead.GetMsgLen() > 0 {
			// 2 再根据 dataLen进行第二次读取, 读出data
			msg := msgHead.(*znet.Message)
			msg.Data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(conn, msg.Data); err != nil {
				fmt.Println("read msg data error: ", err)
				return
			}
			fmt.Println("---> Recv Server Msg: id = ", msg.Id, ", len = ", msg.DataLen, ", data = ", string(msg.Data))
		}

		// cpu阻塞，否则直接爆掉
		time.Sleep(1 * time.Second)
	}
}
