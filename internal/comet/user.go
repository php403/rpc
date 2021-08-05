package comet

import (
	"bufio"
	"github.com/php403/im/api/protocol"
	"sync"
)

type User struct {
	Room     	*Room
	App     	*App
	msgs		chan *protocol.Proto
	cliMsgs		chan *protocol.Proto
	Writer   	*bufio.Writer
	Reader   	*bufio.Reader
	Next     	*User
	Prev     	*User

	AppId    string
	RoomId	 string
	Uid      int64
	token    string
	Server   string
	watchOps map[int32]struct{}
	mutex    sync.RWMutex
}

func NewUser(cli,svr int) *User {
	c := new(User)
	c.cliMsgs = make(chan *protocol.Proto, cli)
	c.msgs = make(chan *protocol.Proto, svr)
	c.watchOps = make(map[int32]struct{})
	return c
}

func (u *User) Push(msg *protocol.Proto) (err error) {
	select {
		case u.msgs <- msg:
	default:
	}
	return
}

func (u *User) Close() {

}

func (u *User) GetMsgs() <-chan *protocol.Proto {
	return u.msgs
}
