package bur

import (
	"sync"

	"log"
)

func notYetImpl(k string, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Printf("%s: not yet implemented", k)
}

func handleError(k string, err error) {
	if err != nil {
		log.Printf("%s: %s", k, err)
	}
}

func Serve(config *Config) {
	log.Printf("AUTH: %s", config.Auth)
	if err := initAuth(config); err != nil {
		handleError("AUTH", err)
		return
	}
	var wg sync.WaitGroup
	for k, c := range config.Proxy {
		switch k {
		case "http":
			wg.Add(1)
			log.Printf("HTTP: %s", c.Addr)
			go func() {
				if err := httpServer(config, &wg); err != nil {
					log.Printf("HTTP: %s", err)
				}
			}()
		case "https":
			wg.Add(1)
			log.Printf("HTTPS: %s %v", c.Addr, c.Params)
			go func() {
				if err := httpsServer(config, &wg); err != nil {
					log.Printf("HTTPS: %s", err)
				}
			}()
		case "socks4":
			wg.Add(1)
			log.Printf("SOCKS4: %s", c.Addr)
			go notYetImpl(k, &wg)
		case "socks5":
			wg.Add(1)
			log.Printf("SOCKS5: %s", c.Addr)
			go notYetImpl(k, &wg)
		case "vpn":
			wg.Add(1)
			log.Printf("VPN: %s", c.Addr)
			go notYetImpl(k, &wg)
		}
	}
	wg.Add(1)
	go func() {
		if err := stateServer(config, &wg); err != nil {
			log.Printf("STATE: %s", err)
		}
	}()
	wg.Wait()
}
