package ziface


/*
	IRequest接口：实际上是把客户端请求的链接信息和请求的数据包装到了一个Request中
*/

type IRequest interface {
	GetConnection() IConnection // 得到当前的链接
	GetData() []byte // 得到请求的消息数据
	GetMsgId() uint32 // 得到消息ID
}