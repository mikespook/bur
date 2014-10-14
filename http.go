package bur

import (
	"log"
	"net/http"
	"sync"

	"github.com/elazarl/goproxy"
	"github.com/elazarl/goproxy/ext/auth"
)

func httpServer(wg *sync.WaitGroup) {
	defer wg.Done()
	proxy := goproxy.NewProxyHttpServer()
	if _config.Auth != "" {
		auth.ProxyBasic(proxy, "bur", authHandle)
	}
	if _config.Debug {
		proxy.Verbose = true
	}
	sc := _config.Proxy["http"]
	if err := http.ListenAndServe(sc.Addr, proxy); err != nil {
		log.Println(err)
	}
}

func httpsServer(wg *sync.WaitGroup) {
	defer wg.Done()
	proxy := goproxy.NewProxyHttpServer()
	if _config.Auth != "" {
		auth.ProxyBasic(proxy, "bur", authHandle)
	}
	if _config.Debug {
		proxy.Verbose = true
	}
	sc := _config.Proxy["https"]
	certFile := sc.Params["cert"]
	keyFile := sc.Params["key"]
	if err := http.ListenAndServeTLS(sc.Addr, certFile, keyFile, proxy); err != nil {
		log.Println(err)
	}
}
