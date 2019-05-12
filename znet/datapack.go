package znet

import (
	"bytes"
	"encoding/binary"
	"github.com/pkg/errors"
	"zinx/utils"
	"zinx/ziface"
)

//封包拆包类实例，暂时不需要成员
type DataPack struct{}

//封包拆包初始化方法
func NewDataPack() *DataPack {
	return &DataPack{}
}

//获取包头长度方法,包头长度为固定,注意再强调一遍
func (dp *DataPack) GetHeadLen() uint32 {
	//Id uint32（4字节） + DataLen uint32(4字节)
	return 8
}

//封包方法
func (dp *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	//创建一个存放bytes字节的缓冲
	dataBuff := bytes.NewBuffer([]byte{})

	/***********************这里写的是包头，写两个int，一共8字节*************************/
	//写dataLen   uint 4字节
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetDataLen()); err != nil {
		return nil, err
	}

	//写msgID    uint 4字节
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}
	/*******************************************************************************/

	/********************************这里写的是data数据********************************/
	//写data数据
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}

	return dataBuff.Bytes(), nil
}

//拆包方法（解压数据）
func (dp *DataPack) Unpack(binaryData []byte) (ziface.IMessage, error) {
	//创建一个从输入二进制数据的ioReader
	dataBuff := bytes.NewReader(binaryData)

	//只解压head的信息，得到dataLen和msgID
	msg := &Message{}

	// 读取dataLen
	// 第三个参数就是要传个变量的指针或者slice，
	// 然后Read函数会按照这个变量类型从buff中读取相应的大小，并按照这变量类型做好相应的字节顺序，
	// 最后把结果赋值给传进来的指针 *data = order.Uint32(bs)
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}

	//读取msgId
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}

	//判断dataLen的长度是否超出我们允许的最大包长度
	if (utils.G_Obj.MaxPacketSize > 0 && msg.DataLen > utils.G_Obj.MaxPacketSize) {
		return nil, errors.New("Too large msg data received")
	}

	//这里只需要把head的数据拆包出来就可以了，然后再通过head的长度，再从conn读取一次数据
	return msg, nil
}
