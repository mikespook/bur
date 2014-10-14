package bur

import (
	"net/http"
	"sync"

	"github.com/elazarl/goproxy"
	"github.com/elazarl/goproxy/ext/auth"
)

func newProxy() (proxy *goproxy.ProxyHttpServer) {
	proxy = goproxy.NewProxyHttpServer()
	if _config.Auth != "" {
		auth.ProxyBasic(proxy, "bur", authHandle)
	}
	if _config.Debug {
		proxy.Verbose = true
	}
	return proxy
}

func httpServer(wg *sync.WaitGroup) error {
	defer wg.Done()
	proxy := newProxy()
	sc := _config.Proxy["http"]
	return http.ListenAndServe(sc.Addr, proxy)
}

func httpsServer(wg *sync.WaitGroup) error {
	defer wg.Done()
	proxy := newProxy()
	sc := _config.Proxy["https"]
	certFile := sc.Params["cert"]
	keyFile := sc.Params["key"]
	return http.ListenAndServeTLS(sc.Addr, certFile, keyFile, proxy)
}
