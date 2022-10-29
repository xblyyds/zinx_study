package znet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"go-study/zinx/utils"
	"go-study/zinx/ziface"
)

// 封包, 拆包的具体模块
type DataPack struct{}

// 创建一个DataPack
func NewDataPack() *DataPack {
	return &DataPack{}
}

func (dp *DataPack) GetHeadLen() uint32 {
	// DataLen uint32 (4字节) + ID uint32 (4字节) = 8字节  这个可以自己定义
	return 8
}

// Pack封包
// |dataLen|msgId|data
func (dp *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	// 创建一个bytes字节的缓冲
	dataBuff := bytes.NewBuffer([]byte{})

	// 协议定义用 小字端 LittleEndian

	// 将dataLen写入dataBuff中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgLen()); err != nil {
		return nil, err
	}

	// 将msgId写入dataBuff中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}

	// 将data数据写入dataBuff中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}

	return dataBuff.Bytes(), nil
}

// Unpack 拆包，将包的Head信息读出来，再根据Head信息里的data长度，再进行一次读
func (dp *DataPack) UnPack(binaryData []byte) (ziface.IMessage, error) {
	// 创建一个从输入二进制数据的ioReader
	dataBuff := bytes.NewReader(binaryData)
	// 只解压head信息, 得到dataLen和msgId
	msg := &Message{}

	// 读dataLen
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}

	// 读msgId
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}

	// 判断dataLen是否已经超出允许的最大包长度 MaxPackageSize
	if utils.GlobalObject.MaxPackageSize > 0 && msg.DataLen > utils.GlobalObject.MaxPackageSize {
		return nil, errors.New("too large msg data recv!")
	}

	return msg, nil
}
