package ziface

/**
 将请求的消息封装到message中，定义抽象接口
 */

type IMessage interface {
	GetDataLen() uint32 //获取消息数据长度
	GetMsgId() uint32   //获取消息ID
	GetData() []byte    //获取消息内容

	SetMsgId(uint32)   //设计消息ID
	SetData([]byte)    //设计消息内容
	SetDataLen(uint32) //设置消息长度
}
