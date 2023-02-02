package znet

import "github.com/hjy497373150/My_zinx/ziface"

// 实现router时，先嵌入这个BaseRouter基类，然后根据需求对基类方法进行重写即可
type BaseRouter struct {

}

// 之所以BaseRouter的方法都为空，是因为有的Router不希望有Prehandle或Posthandle这两个业务
// Router 全部继承Baserouter的好处是不需要实现每一个方法

func (br *BaseRouter) PreHandle(request ziface.IRequest) {

}

func (br *BaseRouter) Handle(request ziface.IRequest) {
	
}

func (br *BaseRouter) PostHandle(request ziface.IRequest) {
	
}