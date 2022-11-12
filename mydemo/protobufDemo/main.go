package main

import (
	"fmt"
	"go-study/mydemo/protobufDemo/pb"
	"google.golang.org/protobuf/proto"
)

func main() {
	person := &pb.Person{
		Name:   "但求一败",
		Age:    16,
		Emails: []string{"https://qq2368269411@163.com", "https://2368269411@qq.com"},
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

	data, err := proto.Marshal(person)
	if err != nil {
		fmt.Println("marshal err:", err)
	}

	newdata := &pb.Person{}
	err = proto.Unmarshal(data, newdata)
	if err != nil {
		fmt.Println("unmarshal err:", err)
	}

	fmt.Println(newdata)
}
