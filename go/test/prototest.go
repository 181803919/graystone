package main

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"os"
	"prototest"
)

func main() {
	fmt.Println("hello world")
	msgTest := &prototest.Person{
		Name: proto.String("stone"),
		Age: proto.Int(33),
		From: proto.String("hubei"),
	}

	msgDataEncode, err := proto.Marshal(msgTest)
	if err != nil{
		panic(err.Error())
		return
	}

	msgEntity := prototest.Person{}
	err = proto.Unmarshal(msgDataEncode, &msgEntity)
	if err != nil{
		fmt.Println(err.Error())
		os.Exit(1)
		return
	}

	fmt.Println("姓名:", msgEntity.GetName())
	fmt.Println("年龄:", msgEntity.GetAge())
	fmt.Println("籍贯:", msgEntity.GetFrom())
}
