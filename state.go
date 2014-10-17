package bur

import (
	"errors"
	"net"
	"net/rpc"
	"os"
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

type StateServer struct {
	l                net.Listener
	network, address string
}

func (srv *StateServer) Close() {
	if srv.l != nil {
		srv.l.Close()
	}
	if srv.network == "unix" || srv.network == "unixpacket" {
		os.Remove(srv.address)
	}
}

func (srv *StateServer) Serve() {
	rpc.Accept(srv.l)
}

var stateServer *StateServer

func newStateServer(config *Config) (err error) {
	srv := &StateServer{}
	defer func() {
		if err == nil {
			stateServer = srv
		}
	}()
	srv.network, srv.address, err = parseNewAddr(config.State.Addr)
	if err != nil {
		return
	}
	bur := &Bur{}
	if err = rpc.RegisterName("Bur", bur); err != nil {
		return
	}
	srv.l, err = net.Listen(srv.network, srv.address)
	return
}

func serveStateServer() {
	stateServer.Serve()
}

func closeStateServer() {
	stateServer.Close()
}
