package bur

import (
	"log"
	"net/rpc"
	"net/url"
	"strings"
)

const (
	Plain     = "plain://"
	Anonymous = "anonymous"
)

var (
	defaultAuth Auth = &authFobidden{}
	users       map[string]string
)

type Auth interface {
	Login(username string, password string) bool
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

func (auth *authRPC) Login(username string, password string) bool {
	var success bool
	if err := auth.Call("Bur.Login",
		map[string]string{"UserName": username, "Password": password},
		&success); err != nil {
		log.Println(err)
		return false
	}
	return success
}

type authPlain struct {
	UserName string
	Password string
}

func (auth *authPlain) Login(username string, password string) bool {
	return auth.UserName == username && auth.Password == password
}

type authFobidden struct{}

func (auth *authFobidden) Login(username string, password string) bool {
	return false
}

type authAnonymous struct{}

func (auth *authAnonymous) Login(username string, password string) bool {
	return true
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
	return &authPlain{username, password}
}

func useRPCAuth(config string) *authRPC {
	u, err := url.Parse(config)
	if err != nil {
		log.Println(err)
		return nil
	}
	auth, err := NewAuthRPC(u.Scheme, u.Host)
	if err != nil {
		log.Println(err)
		return nil
	}
	return auth
}

func initAuth() {
	switch {
	case _config.Auth == Anonymous:
		defaultAuth = &authAnonymous{}
	case strings.HasPrefix(_config.Auth, Plain):
		defaultAuth = usePlainAuth(_config.Auth)
	default:
		defaultAuth = useRPCAuth(_config.Auth)
	}
}

func authHandle(username, password string) bool {
	if password, ok := users[username]; ok {
		return users[username] == password
	}
	if defaultAuth.Login(username, password) {
		users[username] = password
	}
	return false
}
