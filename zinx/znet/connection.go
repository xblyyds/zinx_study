package znet

import (
	"errors"
	"fmt"
	"go-study/zinx/utils"
	"go-study/zinx/ziface"
	"io"
	"net"
)

/*
	链接模块
*/
type Connection struct {
	// 当前链接的socket TCP套接字
	Conn *net.TCPConn

	// 链接的ID
	ConnID uint32

	// 当前的链接状态
	isClosed bool

	// 告知当前链接已经退出的/停止 channel
	ExitChan chan bool

	// 无缓冲的管道, 用于读, 写Groutine之间的通信
	msgChan chan []byte

	// 消息的管理 MsgId 和 对应处理业务API 关系
	MsgHandler ziface.IMsgHandle
}

// 初始化链接模块的方法
func NewConnection(conn *net.TCPConn, connID uint32, msgHandler ziface.IMsgHandle) *Connection {
	c := &Connection{
		Conn:       conn,
		ConnID:     connID,
		isClosed:   false,
		MsgHandler: msgHandler,
		ExitChan:   make(chan bool, 1),
		msgChan:    make(chan []byte),
	}

	return c
}

func (c *Connection) StartReader() {
	fmt.Println("[Reader Goroutine is running]")
	defer fmt.Println("[Reader is exit]", "connID = ", c.ConnID, " , remote addr is ", c.RemoteAddr().String())
	defer c.Stop()

	for {
		// 读取客户端的数据到buf中, 最大512字节
		//buf := make([]byte, utils.GlobalObject.MaxPackageSize)
		//_, err := c.Conn.Read(buf)
		//if err != nil {
		//	fmt.Println("recv buf err", err)
		//	break
		//}

		// 创建一个拆包解包对象
		dp := NewDataPack()
		// 读取客户端的 Msg Head 二进制流 8 个字节
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("read msg head error: ", err)
			break
		}

		// 拆包, 得到 msgId 和 msgDataLen 放到 msg消息中
		msg, err := dp.UnPack(headData)
		if err != nil {
			fmt.Println("unpack error: ", err)
			break
		}

		// 根据 dataLen 再次读取Data, 放在 msg.Data 中
		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("read msg data error: ", err)
				break
			}
		}
		msg.SetData(data)

		// 得到当前conn数据的Request请求数据
		req := Request{
			conn: c,
			msg:  msg,
		}

		if utils.GlobalObject.WorkerPoolSize > 0 {
			c.MsgHandler.SendMsgToTaskQueue(&req)
		} else {
			// 根据绑定好的MsgId找到处理对应API业务 执行
			go c.MsgHandler.DoMsgHandler(&req)

			// 从路由中，找到注册绑定的Conn对应的router调用
		}
	}
}

func (c *Connection) StartWriter() {
	fmt.Println("[Write Goroutine is running]")
	defer fmt.Println(" [conn Writer exit!]", c.RemoteAddr().String())
	// 不断的阻塞等待channel消息, 进行写给客户端
	for {
		select {
		case data := <-c.msgChan:
			// 有数据要写给客户端
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send data error: ", err)
				return
			}
		case <-c.ExitChan:
			// Reader 已经退出, 此时 Writer 也要退出
			return
		}
	}
}

// 启动链接 让当前的链接准备开始工作
func (c *Connection) Start() {
	fmt.Println("Conn Start().. ConnID = ", c.ConnID)
	// 启动从当前链接的读数据的业务
	go c.StartReader()
	// 启动从当前链接写数据的业务
	go c.StartWriter()
}

// 停止链接 结束当前链接的工作
func (c *Connection) Stop() {
	fmt.Println("Conn Stop().. ConnID = ", c.ConnID)

	// 如果当前链接已经关闭
	if c.isClosed == true {
		return
	}

	c.isClosed = true

	// 关闭socket链接
	c.Conn.Close()
	// 告知 Writer 关闭
	c.ExitChan <- true

	close(c.ExitChan)
	close(c.msgChan)
}

// 获取当前链接的绑定socket conn
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

// 获取当前链接模块的链接ID
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

// 获取远程客户端的 TCP的状态 IP Port
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// 发送数据 将数据发送给远程的客户端
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("connection closed when send msg")
	}
	// 将 data进行封包
	dp := NewDataPack()

	binaryMsg, err := dp.Pack(NewMessage(msgId, data))
	if err != nil {
		fmt.Println("pack error msg is = ", msgId)
		return errors.New("pack err msg")
	}

	// 读写分离：将数据发送给消息管道
	c.msgChan <- binaryMsg
	return nil
}
