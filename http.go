package bur

import (
	"net/http"
	"sync"

	"github.com/elazarl/goproxy"
	"github.com/elazarl/goproxy/ext/auth"
)

func newProxy(ads string, debug bool) (proxy *goproxy.ProxyHttpServer) {
	proxy = goproxy.NewProxyHttpServer()
	if ads != "" {
		auth.ProxyBasic(proxy, "bur", authHandle)
	}
	if debug {
		proxy.Verbose = true
	}
	return proxy
}

func httpServer(config *Config, wg *sync.WaitGroup) error {
	defer wg.Done()
	proxy := newProxy(config.Auth, config.Debug)
	sc := config.Proxy["http"]
	return http.ListenAndServe(sc.Addr, proxy)
}

func httpsServer(config *Config, wg *sync.WaitGroup) error {
	defer wg.Done()
	proxy := newProxy(config.Auth, config.Debug)
	sc := config.Proxy["https"]
	certFile := sc.Params["cert"]
	keyFile := sc.Params["key"]
	return http.ListenAndServeTLS(sc.Addr, certFile, keyFile, proxy)
}
