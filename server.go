package bur

import (
	"strings"
	"sync"

	"github.com/mikespook/golib/log"
)

func notYetImpl(k string, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Errorf("%s: not yet implemented", strings.ToUpper(k))
}

func handleError(k string, err error) {
	if err != nil {
		log.Errorf("%s: %s", k, err)
	}
}

func Serve(config *Config) {
	if err := initAuth(config); err != nil {
		handleError("AUTH", err)
		return
	}
	var wg sync.WaitGroup
	for k := range config.Proxy {
		switch k {
		case "http":
			wg.Add(1)
			go handleError("HTTP", httpServer(config, &wg))
		case "https":
			wg.Add(1)
			go handleError("HTTPS", httpsServer(config, &wg))
		case "socks4":
			wg.Add(1)
			go notYetImpl(k, &wg)
		case "socks5":
			wg.Add(1)
			go notYetImpl(k, &wg)
		case "vpn":
			wg.Add(1)
			go notYetImpl(k, &wg)
		}
	}
	wg.Add(1)
	go handleError("STATE", stateServer(config, &wg))
	wg.Wait()
}
