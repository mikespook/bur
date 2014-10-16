package bur

import (
	"sync"

	"log"
)

func notYetImpl(wg *sync.WaitGroup) {
	defer wg.Done()
	log.Printf("not yet implemented")
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
		log.Println(err)
		return false
	} else if permit {
		users.Set(username, password)
		return true
	}
	return false
}

func Serve(config *Config) {
	log.Println(config)
	if err := initAuth(config); err != nil {
		log.Println(err)
		return
	}
	var wg sync.WaitGroup
	for k := range config.Proxy {
		switch k {
		case "http":
			wg.Add(1)
			go func() {
				if err := httpServer(config, &wg); err != nil {
					log.Println(err)
				}
			}()
		case "https":
			wg.Add(1)
			go func() {
				if err := httpsServer(config, &wg); err != nil {
					log.Println(err)
				}
			}()
		case "socks4":
			wg.Add(1)
			go notYetImpl(&wg)
		case "socks5":
			wg.Add(1)
			go notYetImpl(&wg)
		case "vpn":
			wg.Add(1)
			go notYetImpl(&wg)
		}
	}
	wg.Add(1)
	go func() {
		if err := stateServer(config, &wg); err != nil {
			log.Println(err)
		}
	}()
	wg.Wait()
}
