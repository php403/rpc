package comet

import (
	"github.com/php403/im/api/protocol"
	"sync"
)

type Room struct {
	ID        string
	rLock     sync.RWMutex
	next      *User
}

func NewRoom(id string) (r *Room) {
	r = new(Room)
	r.ID = id
	r.next = nil
	return
}

func(room *Room) Put(user *User)  {
	room.rLock.Lock()
	if room.next != nil {
		room.next.Prev = user
	}
	user.Next = room.next
	user.Prev = nil
	room.next = user // insert to header
}

func(room *Room) Del(user *User) bool {
	room.rLock.Lock()
	if user.Prev == nil && user.Next == nil {
		room.rLock.Unlock()
		return false
	}

	if user.Next != nil {
		user.Next.Prev = user.Prev
	}

	if user.Prev != nil {
		user.Next.Next = user.Next
	}else{
		room.next = user.Next
	}
	user.Prev = nil
	user.Next = nil
	room.rLock.Unlock()
	return true
}

func(room *Room) Push(msg *protocol.Proto)  {
	room.rLock.RLock()
	for user := room.next; user != nil; user = user.Next {
		_ = user.Push(msg)
	}
	room.rLock.RUnlock()
}

func (room *Room) Close() {
	room.rLock.RLock()
	for user := room.next; user != nil; user = user.Next {
		user.Close()
	}
	room.rLock.RUnlock()
}