package bur

import (
	"sync"
	"sync/atomic"
)

type State struct {
	Tx, Rx uint64
}

func (state *State) ToRx(delta uint64) {
	atomic.AddUint64(&state.Rx, delta)
}

func (state *State) ToTx(delta uint64) {
	atomic.AddUint64(&state.Tx, delta)
}

type User struct {
	Name, Password string
	State          UserState
}

func (user *User) Logined() {
	atomic.AddUint32(&user.State.Login, 1)
}

type UserState struct {
	Login uint32
	State
}

type Users struct {
	sync.RWMutex
	List map[string]*User
}

func NewUsers() (users Users) {
	users = Users{List: make(map[string]*User)}
	return users
}

func (users *Users) Get(name string) *User {
	defer users.RUnlock()
	users.RLock()
	if user, ok := users.List[name]; ok {
		return user
	}
	return nil
}

func (users *Users) Set(name, password string) {
	defer users.Unlock()
	users.Lock()
	users.List[name] = &User{name, password, UserState{Login: 1}}
}
