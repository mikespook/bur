package bur

import (
	"log"
	"sync"
)

func notYetImpl(k string, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Println(k, ": not yet implemented")
}

func Serve() {
	initAuth()
	var wg sync.WaitGroup
	for k := range _config.Proxy {
		switch k {
		case "http":
			wg.Add(1)
			go httpServer(&wg)
		case "https":
			wg.Add(1)
			go httpsServer(&wg)
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
	go ctrlServer(&wg)
	wg.Wait()
}
