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
	users            = make(map[string]string)
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

func (auth *authRPC) Login(username string, password string) (permit bool, err error) {
	if err := auth.Call("Bur.Login",
		map[string]string{"UserName": username, "Password": password},
		&permit); err != nil {
		return false, err
	}
	return
}

type authPlain struct {
	UserName string
	Password string
}

func (auth *authPlain) Login(username string, password string) (bool, error) {
	return auth.UserName == username && auth.Password == password, nil
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
	if permit, err := defaultAuth.Login(username, password); err != nil {
		handleError("AUTH", err)
		return false
	} else if permit {
		users[username] = password
		return true
	}
	return false
}
