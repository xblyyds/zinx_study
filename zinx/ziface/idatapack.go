package ziface

// 封包, 拆包抽象接口
// 直接面向TCP链接中的数据留, 用于处理TCP粘包问题

type IDataPack interface {
	// 获取包头的长度
	GetHeadLen() uint32

	// 封包
	Pack(IMessage) ([]byte, error)

	// 拆包
	UnPack([]byte) (IMessage, error)
}
