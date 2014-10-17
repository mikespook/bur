package bur

import (
	"net/url"
	"strings"
	"sync"

	"github.com/mikespook/golib/log"
)

func notYetImpl(k string, wg *sync.WaitGroup) {

}

func authHandle(username, password string) (ok bool) {
	defer func() {
		if ok {
			users.Get(username).Logined()
		}
	}()
	if user := users.Get(username); user != nil && user.password == password {
		return true
	}
	if permit, err := defaultAuth.Login(username, password); err != nil {
		log.Error(err)
		return false
	} else if permit {
		users.Set(username, password)
		return true
	}
	return false
}

func Serve(config *Config) {
	log.Messagef("AUTH: %v", config.Auth)
	if err := initAuth(config); err != nil {
		log.Error(err)
		return
	}
	for k := range config.Proxy {
		switch k {
		case "http":
			go func() {
				log.Messagef("HTTP: %v", config.Proxy["http"])
				if err := httpServer(config); err != nil {
					log.Error(err)
				}
			}()
		case "https":
			go func() {
				log.Messagef("HTTPS: %v", config.Proxy["https"])
				if err := httpsServer(config); err != nil {
					log.Error(err)
				}
			}()
		case "socks4":
			log.Warning("SOCKS4: not yet implemented")
		case "socks5":
			log.Warning("SOCKS5: not yet implemented")
		case "vpn":
			log.Warning("VPN: not yet implemented")
		}
	}
	if err := newStateServer(config); err != nil {
		log.Error(err)
	}
	go serveStateServer()
}

func Close() {
	closeStateServer()
}

func parseNewAddr(ads string) (network, address string, err error) {
	ads = strings.Replace(ads, ":///", "://file/", 1)
	u, err := url.Parse(ads)
	if err != nil {
		return
	}
	network = u.Scheme
	if u.Host == "file" {
		address = u.Path
	} else {
		address = u.Host
	}
	return
}
