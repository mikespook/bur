package bur

import (
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
	if err := a.Call("Bur.Login",
		map[string]string{"UserName": username, "Password": password},
		&permit); err != nil {
		return false, err
	}
	return
}

type authPlain struct {
	User
}

func (auth *authPlain) Login(username string, password string) (bool, error) {
	return auth.Name == username && auth.Password == password, nil
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
	u, err := url.Parse(config)
	if err != nil {
		return
	}
	return NewAuthRPC(u.Scheme, u.Host)
}

func initAuth(config *Config) (err error) {
	switch {
	case config.Auth == Anonymous:
		defaultAuth = &authAnonymous{}
	case strings.HasPrefix(config.Auth, Plain):
		defaultAuth = usePlainAuth(config.Auth)
	default:
		defaultAuth, err = useRPCAuth(config.Auth)
	}
	return
}

func cachedAuth(username, password string) bool {
	defer usersLock.RUnlock()
	usersLock.RLock()
	if user, ok := users[username]; ok {
		usersLock.RUnlock()
		match := (user.Password == password)
		user.State.Login++
		return match
	}
	return false
}

func authHandle(username, password string) bool {
	if ok := cachedAuth(username, password); ok {
		return true
	}
	if permit, err := defaultAuth.Login(username, password); err != nil {
		handleError("AUTH", err)
		return false
	} else if permit {
		defer usersLock.Unlock()
		usersLock.Lock()
		users[username] = User{username, password, UserState{Login: 1}}
		return true
	}
	return false
}
