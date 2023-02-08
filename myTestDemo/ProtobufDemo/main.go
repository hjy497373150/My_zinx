package main

import (
	"My_zinx/myTestDemo/ProtobufDemo/pb"
	"fmt"

	"github.com/golang/protobuf/proto"
)

func main() {
	person := &pb.Person{
		Name: "klayhu",
		Age: 23,
		Emails: "493737150@qq.com",
		Phones: []*pb.PhoneNumber{
            &pb.PhoneNumber{
                Number: "13113111311",
                Type:   pb.PhoneType_MOBILE,
            },
            &pb.PhoneNumber{
                Number: "14141444144",
                Type:   pb.PhoneType_HOME,
            },
            &pb.PhoneNumber{
                Number: "19191919191",
                Type:   pb.PhoneType_WORK,
            },
        },
	}
	// 编码 将person对象进行序列化，得到二进制文件
	data ,err := proto.Marshal(person)
	// data就是我们要进行网络传输的数据，对端需要按照message person格式进行解析
	if err != nil {
		fmt.Println("marshal error: ",err)
	}

	// 解码
	newdata := &pb.Person{}
	err = proto.Unmarshal(data, newdata)

	if err != nil {
		fmt.Println("unmarshal error: ",err)
	}

	fmt.Println("源数据:",data)
	fmt.Println("解码之后的数据",newdata)

}