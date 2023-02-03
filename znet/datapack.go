package znet

import (
	"bytes"
	"encoding/binary"
	"errors"

	"github.com/hjy497373150/My_zinx/utils"
	"github.com/hjy497373150/My_zinx/ziface"
)

type DataPack struct {

}

// 获取包的头部长度方法
func (d *DataPack) GetHeadLen() uint32 {
	// message中，datalen是uint32 4字节，Id也是uint32 4字节
	return 8
}

// 封包方法
// datalen|Id|data
func (d *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	// 创建一个存放byte字节的缓冲
	databuf := bytes.NewBuffer([]byte{})

	// 将datalen写到databuf中
	if err := binary.Write(databuf, binary.LittleEndian, msg.GetMsgLen());err != nil {
		return nil,err
	}

	// 将ID写到databuf中
	if err := binary.Write(databuf, binary.LittleEndian, msg.GetMsgId());err != nil {
		return nil,err
	}

	//将data写到databuf中
	if err := binary.Write(databuf, binary.LittleEndian, msg.GetData());err != nil {
		return nil,err
	}

	return databuf.Bytes(),nil
}

// 拆包方法,将包的head信息读出来，根据head信息里data的长度，再进行一次读
func (d *DataPack) UnPack(binaryData []byte) (ziface.IMessage,error) {
	// 创建一个从输入二进制数据的ioreader
	databuf := bytes.NewReader(binaryData)

	// 只解压head的信息，得到datalen和id
	msg := &Message{}
	//读datalen
	if err := binary.Read(databuf,binary.LittleEndian,&msg.DataLen);err != nil {
		return nil,err
	}
	//读id
	if err := binary.Read(databuf,binary.LittleEndian,&msg.ID);err != nil {
		return nil,err
	}

	// 判断datalen的长度是否超过Server允许包体的最大长度
	if utils.GlobalObject.MaxPackageSize > 0 && msg.DataLen > utils.GlobalObject.MaxPackageSize {
		return nil,errors.New("Too large msg data recv...")
	}

	return msg,nil

}

// 拆包封包的初始化方法
func NewDataPack() *DataPack {
	return &DataPack{}
}