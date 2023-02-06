package utils

import (
	"encoding/json"
	"io/ioutil"

	"github.com/hjy497373150/My_zinx/ziface"
)

/*
	存储一切有关Zinx框架的全局参数，供其他模块使用
	一些参数是可以通过zinx.json由用户进行配置

*/

type GlobalObj struct {
	/*
		Server
	*/
	TcpServer ziface.IServer //当前Zinx全局的Server对象
	Host string // 当前服务器主机监听的ip
	TcpPort int // 当前服务器主机监听的端口号
	Name string // 当前服务器的名称

	/*
		Myzinx
	*/
	Version string // 当前myzinx的版本号
	MaxConn int // 当前服务器主机允许的最大链接数
	MaxPackageSize uint32 // 当前Zinx框架数据包的最大值
	WorkerPoolSize uint32 //业务工作Worker池的数量
	MaxWorkTaskLen uint32 //每一个业务工作Worker池对应的消息队列的最大长度

	/*
		config file path
	*/
	ConfFilePath string
}

// 从myzinx.json去加载用于自定义的参数
func (g *GlobalObj) Reload() {
	data,err := ioutil.ReadFile("conf/myzinx.json")
	if err != nil {
		panic(err) 
	}

	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}


// 定义一个全局的对外GlobalObj
var GlobalObject *GlobalObj

// 提供init方法，初始化当前的GlobalObject
func init() {
	// 如果配置文件没有加载，默认值
	GlobalObject := &GlobalObj{
		Name: "MyZinxServer",
		Version: "V0.8",
		TcpPort: 8888,
		Host: "127.0.0.1",
		MaxConn: 1000,
		MaxPackageSize: 4096,
		WorkerPoolSize: 10,
		MaxWorkTaskLen: 1024,
		ConfFilePath:  "conf/myzinx.json",
	}

	GlobalObject.Reload()
}
