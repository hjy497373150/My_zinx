package znet

import (
	"fmt"
	"strconv"

	"github.com/hjy497373150/My_zinx/ziface"
)

type MsgHandle struct {
	Apis map[uint32] ziface.IRouter // 存放每个msgid所对应的处理方法的map
}

// 以非阻塞方式处理消息
func (mh *MsgHandle) DoMsgHandle(request ziface.IRequest) {
	handle,ok := mh.Apis[request.GetMsgId()];
	if !ok {
		fmt.Println("api msgId:",request.GetMsgId(), " not Found!")
		return
	}
	// 执行对应的处理方法
	handle.Handle(request)
}


// 为消息添加具体的处理逻辑
func (mh *MsgHandle) AddRouter(msgId uint32, router ziface.IRouter) {
	// 1.判断当前msgid所绑定的api处理方法是否已经存在
	if _,ok := mh.Apis[msgId];ok {
		panic("Repeat api, msgId = " + strconv.Itoa(int(msgId)))
	}
	// 2.添加新的绑定关系
	mh.Apis[msgId] = router
	fmt.Println("Add api msgId:",msgId)
}

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis: make(map[uint32] ziface.IRouter),
	}
}