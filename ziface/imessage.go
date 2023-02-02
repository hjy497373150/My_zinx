package ziface

/*
	将请求的消息封装到一个Message中，定义抽象的接口
*/
type IMessage interface {
	GetMsgId() uint32 // 获取消息的长度
	GetMsgLen() uint32 // 获取消息的ID
	GetData() []byte // 获取消息内容

	SetMsgId(uint32)  // 设置消息的长度
	SetMsgLen(uint32)  // 设置消息的ID
	SetData([]byte) // 设置消息的内容

}