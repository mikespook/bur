package bur

import (
	"errors"
	"net"
	"net/rpc"
	"net/url"
	"sync"
)

var users = NewUsers()
var total State
var errUserNotFound = errors.New("User Not Found")

type Bur struct{}

func (b *Bur) User(name string, user *User) error {
	if *user = *users.Get(name); user == nil {
		return errUserNotFound
	}
	return nil
}

func (b *Bur) State(a interface{}, state *State) error {
	*state = total
	return nil
}

func stateServer(config *Config, wg *sync.WaitGroup) error {
	defer wg.Done()
	bur := &Bur{}
	rpc.Register(bur)
	u, err := url.Parse(config.State)
	if err != nil {
		return err
	}
	l, err := net.Listen(u.Scheme, u.Host)
	if err != nil {
		return err
	}
	rpc.Accept(l)
	return nil
}
