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
	u := users.Get(name)
	if u == nil {
		return errUserNotFound
	}
	*user = *u
	return nil
}

func (b *Bur) State(a string, state *State) error {
	*state = total
	return nil
}

func stateServer(config *Config, wg *sync.WaitGroup) error {
	defer wg.Done()
	bur := &Bur{}
	if err := rpc.RegisterName("Bur", bur); err != nil {
		return err
	}
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
