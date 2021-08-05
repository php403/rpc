package main

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/php403/im/api/logic"
	"github.com/php403/im/api/protocol"
)

func main()  {
	data := []byte("hello")
	pb := &logic.PushMsg{
		Msg: data,
		AppId: "1",
		RoomId: "1",
		UserId: 1,
		Server:"1",
		Operation:protocol.OpSendAppMsg,
	}
	pb.Type = logic.PushMsg_PUSH
	msg,err := proto.Marshal(pb)
	if err != nil {
		panic(err)
	}
	fmt.Println(msg)
	
}
