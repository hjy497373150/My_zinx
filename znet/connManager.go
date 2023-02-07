package znet

import (
	"errors"
	"fmt"
	"sync"

	"github.com/hjy497373150/My_zinx/ziface"
)

type ConnManager struct {
	Connections map[uint32] ziface.IConnection //管理的链接集合
	connLock sync.RWMutex // 保护链接集合的读写锁
}

// ConnManager的初始化方法
func NewConnManager() *ConnManager {
	return &ConnManager{
		Connections: make(map[uint32]ziface.IConnection),
	}
}

// 添加链接
func (connMgr *ConnManager) Add(conn ziface.IConnection) {
	// 保护共享资源 加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	connMgr.Connections[conn.GetConnID()] = conn

	fmt.Println("conn id = ",conn.GetConnID(),"is added to connMgr successfully,connMgrLen = ",connMgr.Len())
}

// 删除链接
func (connMgr *ConnManager) Remove(conn ziface.IConnection) {
	// 保护共享资源 加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()
	if connection,ok := connMgr.Connections[conn.GetConnID()];ok {
		delete(connMgr.Connections,connection.GetConnID())
		fmt.Println("conn id = ",conn.GetConnID(),"is removed from connMgr successfully,connMgrLen = ",connMgr.Len())
	} else {
		fmt.Println("conn id = ",conn.GetConnID(),"connMgr not found this conn,remove failed!!!")
	}
}

// 根据connId获取链接
func (connMgr *ConnManager) Get(connId uint32) (ziface.IConnection,error) {
	// 保护共享资源 加读锁
	connMgr.connLock.RLock()
	defer connMgr.connLock.RUnlock()

	if connection,ok := connMgr.Connections[connId];ok {
		fmt.Println("conn id = ",connId,"has found in connMgr successfully")
		return connection,nil
		
	} else {
		return nil,errors.New("connMgr not found this conn")
	}
}

// 得到当前链接总数
func (connMgr *ConnManager) Len() int {
	return len(connMgr.Connections)
}

// 清理所有链接
func (connMgr *ConnManager) ClearAll() {
	// 保护共享资源 加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	// 停止并删除所有的链接信息
	for connId,conn := range connMgr.Connections {
		// 停止
		conn.Stop()
		// 删除
		delete(connMgr.Connections,connId)
	}

	fmt.Println("All connections have been cleared successfully!")
}