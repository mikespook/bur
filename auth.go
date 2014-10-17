package bur

import (
	"net/rpc"
	"strings"
)

const (
	Plain     = "plain://"
	Anonymous = "anonymous"
)

var (
	defaultAuth Auth = &authFobidden{}
)

type Auth interface {
	Login(username string, password string) (bool, error)
}

func NewAuthRPC(network, address string) (auth *authRPC, err error) {
	client, err := rpc.Dial(network, address)
	if err != nil {
		return
	}
	auth = &authRPC{client}
	return
}

type authRPC struct {
	*rpc.Client
}

func (a *authRPC) Login(username string, password string) (permit bool, err error) {
	if err := a.Call("User.Login",
		map[string]string{"Name": username, "Password": password},
		&permit); err != nil {
		return false, err
	}
	return
}

type authPlain struct {
	User
}

func (auth *authPlain) Login(username string, password string) (bool, error) {
	return auth.Name == username && auth.password == password, nil
}

type authFobidden struct{}

func (auth *authFobidden) Login(username string, password string) (bool, error) {
	return false, nil
}

type authAnonymous struct{}

func (auth *authAnonymous) Login(username string, password string) (bool, error) {
	return true, nil
}

func usePlainAuth(config string) *authPlain {
	config = strings.TrimPrefix(config, Plain)
	tmp := strings.SplitN(config, ":", 2)
	var username, password string
	if len(tmp) == 2 {
		username = tmp[0]
		password = tmp[1]
	} else {
		username = tmp[0]
		password = ""
	}
	return &authPlain{User{username, password, UserState{}}}
}

func useRPCAuth(config string) (aRPC *authRPC, err error) {
	network, address, err := parseNewAddr(config)
	if err != nil {
		return
	}
	return NewAuthRPC(network, address)
}

func initAuth(config *Config) (err error) {
	switch {
	case config.Auth.Addr == Anonymous:
		defaultAuth = &authAnonymous{}
	case strings.HasPrefix(config.Auth.Addr, Plain):
		defaultAuth = usePlainAuth(config.Auth.Addr)
	default:
		defaultAuth, err = useRPCAuth(config.Auth.Addr)
	}
	return
}
