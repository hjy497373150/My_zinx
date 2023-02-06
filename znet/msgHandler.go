package znet

import (
	"fmt"
	"strconv"

	"github.com/hjy497373150/My_zinx/utils"
	"github.com/hjy497373150/My_zinx/ziface"
)

type MsgHandle struct {
	Apis map[uint32] ziface.IRouter // 存放每个msgid所对应的处理方法的map
	WorkerPoolSize uint32 // 业务工作Worker池的大小
	TaskQueue []chan ziface.IRequest //Worker负责取任务的消息队列
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

// 启动一个Worker工作池（开启工作池的动作只能发生一次）
func (mh *MsgHandle)StartWorkPool() {
	// 根据WorkerPoolSize 分别开启Worker 每个Worker用一个go来承载
	for i:=0;i<int(mh.WorkerPoolSize);i++ {
		// 一个worker被启动
		// 1.当前的worker对应的channel消息队列 开辟空间 第0个worker就用第0个channel
		mh.TaskQueue[i] = make(chan ziface.IRequest,utils.GlobalObject.MaxWorkTaskLen)
		// 2.启动当前的Worker 阻塞等待消息从channel传递进来
		go mh.StartOneWorker(i,mh.TaskQueue[i])

	}
}

// 启动一个Worker工作流程
func (mh *MsgHandle)StartOneWorker(workerId int, taskQueue chan ziface.IRequest) {
	fmt.Println("WorkId = ",workerId, " is started...")

	// 不断地阻塞等待消息队列的消息
	for {
		select {
			// 如果有消息过来，出列的就是一个客户端的Request，执行这个Request对应的handle方法
		case request := <- taskQueue:
			mh.DoMsgHandle(request)
		}
	}
}

// 将消息交给消息处理队列，让Worker处理
func (mh *MsgHandle)SendMsgToTaskQueue(request ziface.IRequest) {
	// 1.将消息平均分配给不同的Worker，基本的轮询
	workerId := request.GetConnection().GetConnID() % mh.WorkerPoolSize

	fmt.Println("Add Conn ID = ",request.GetConnection().GetConnID()," request msgId = ",request.GetMsgId()," to workerId = ",workerId)
	// 2.将消息发送给对应的Worker的TaskQueue即可
	mh.TaskQueue[workerId] <- request

}


func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis: make(map[uint32] ziface.IRouter),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
		TaskQueue: make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize), //一个Worker对应一个Task
	}
}