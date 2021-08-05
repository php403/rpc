package comet

import (
	"github.com/php403/im/api/protocol"
	"sync"
)

type App struct {
	ID        string
	rLock     sync.RWMutex
	next      *User
}

func NewApp(id string) (app *App) {
	app = new(App)
	app.ID = id
	app.next = nil
	return
}

func(app *App) Put(user *User)  {
	app.rLock.Lock()
	if app.next != nil {
		app.next.Prev = user
	}
	user.Next = app.next
	user.Prev = nil
	app.next = user // insert to header
	app.rLock.Unlock()
}

func(app *App) Del(user *User) bool {
	app.rLock.Lock()
	if user.Prev == nil && user.Next == nil {
		app.rLock.Unlock()
		return false
	}

	if user.Next != nil {
		user.Next.Prev = user.Prev
	}

	if user.Prev != nil {
		user.Next.Next = user.Next
	}else{
		app.next = user.Next
	}
	user.Prev = nil
	user.Next = nil
	app.rLock.Unlock()
	return true
}

func(app *App) Push(msg *protocol.Proto)  {
	app.rLock.RLock()
	for user := app.next; user != nil; user = user.Next {
		_ = user.Push(msg)
	}
	app.rLock.RUnlock()
}

func (app *App) Close() {
	app.rLock.RLock()
	for user := app.next; user != nil; user = user.Next {
		user.Close()
	}
	app.rLock.RUnlock()
}

