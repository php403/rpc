package protocol

var Ver = uint32(1)
const (
	OpHeartbeat = uint32(1)
	OpHeartbeatReply = uint32(2)

	OpAuth = uint32(3)
	OpAuthReply = uint32(4)

	OpSendMsg = uint32(5)
	OpSendMsgReply = uint32(6)

	OpSendAppMsg = uint32(7)
	OpSendAppMsgReply = uint32(8)

	OpSendRoomMsg = uint32(9)
	OpSendRoomMsgReply = uint32(10)

	OpClose = uint32(11)
	OpClosetReply = uint32(12)
)




func GetCloseProto() *Proto {
	return &Proto{
		Ver:Ver,
		Op : OpClosetReply,
		Body:[]byte("server close!"),

	}
}
