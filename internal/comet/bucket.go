package comet

import (
	"github.com/php403/im/api/protocol"
	"github.com/php403/im/internal/comet/conf"
	"sync"
)

type Bucket struct {
	c     *conf.Bucket
	cLock sync.RWMutex        // protect the channels for users
	users   map[int64]*User // map sub key to a channel
	// room
	apps       map[string]*App // bucket app channels
	rooms       map[string]*Room // bucket room channels
	//routines    []chan *comet.BroadcastRoomReq//*pb.BroadcastRoomReq
	routinesNum uint64
}

func NewBucket(c *conf.Bucket) (b *Bucket) {
	b = new(Bucket)
	b.users = make(map[int64]*User, c.Channel)
	b.c = c
	b.apps =  make(map[string]*App, c.App)
	b.rooms =  make(map[string]*Room, c.App)
	//b.routines = make([]chan *comet.BroadcastRoomReq, c.RoutineAmount)
	for i := uint64(0); i < c.RoutineAmount; i++ {
		//c := make(chan *comet.BroadcastRoomReq, c.RoutineSize)
		//b.routines[i] = c
	}
	return
}

func(b *Bucket) Put(appId string,roomId string,user *User) (err error) {
	var ok bool
	var room *Room
	var app *App
	b.cLock.Lock()
	if user := b.users[user.Uid]; user != nil {
		user.Close()
	}
	if roomId != "" {
		if room, ok = b.rooms[roomId]; !ok {
			room = NewRoom(roomId)
			b.rooms[roomId] = room
		}
		user.Room = room
		room.Put(user)
	}

	if app,ok = b.apps[appId]; !ok {
		app = NewApp(appId)
		b.apps[appId] = app
	}
	user.App = app
	app.Put(user)
	b.cLock.Unlock()
	return
}

func (b *Bucket) Del(user *User) {
	var (
		ok   	bool
		user1    *User
	)
	b.cLock.Lock()
	if user1, ok = b.users[user.Uid]; ok {
		if user == user1 {
			delete(b.users, user.Uid)
		}

	}
	b.cLock.Unlock()
	//todo if room.users=nil delete *room
}

func (b *Bucket) User(uid int64) (ch *User) {
	b.cLock.RLock()
	ch = b.users[uid]
	b.cLock.RUnlock()
	return
}

func (b *Bucket) Broadcast(p *protocol.Proto, op int32) {
	var user *User
	b.cLock.RLock()
	//todo 可添加广播过滤
	for _, user = range b.users {
		_ = user.Push(p)
	}
	b.cLock.RUnlock()
}

// Room get a room by roomid.
func (b *Bucket) Room(rid string) (room *Room) {
	b.cLock.RLock()
	room = b.rooms[rid]
	b.cLock.RUnlock()
	return
}



/*func (b *Bucket) BroadcastRoom(arg *comet.BroadcastRoomReq) {
	num := atomic.AddUint64(&b.routinesNum, 1) % b.c.RoutineAmount
	b.routines[num] <- arg
}
*/


