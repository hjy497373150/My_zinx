package znet

type Message struct {
	ID uint32 // 消息的id
	DataLen uint32 // 消息的长度
	Data []byte // 消息的内容
}

// 获取消息的长度
func (msg *Message)GetMsgId() uint32 {
	return msg.ID
}

// 获取消息的ID
func (msg *Message) GetMsgLen() uint32 {
	return msg.DataLen
}

// 获取消息内容
func (msg *Message)GetData() []byte {
	return msg.Data
}

// 设置消息的长度
func (msg *Message)SetMsgId(id uint32) {
	msg.ID = id
}

// 设置消息的ID
func (msg *Message)SetMsgLen(len uint32) {
	msg.DataLen = len
}

// 设置消息的内容
func (msg *Message)SetData(data []byte) {
	msg.Data = data
}