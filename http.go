package bur

import (
	"net/http"

	"github.com/elazarl/goproxy"
	"github.com/elazarl/goproxy/ext/auth"
)

func newProxy(ads string) (proxy *goproxy.ProxyHttpServer) {
	proxy = goproxy.NewProxyHttpServer()
	if ads != "" {
		auth.ProxyBasic(proxy, "bur", authHandle)
	}
	return proxy
}

func httpServer(config *Config) error {
	proxy := newProxy(config.Auth.Addr)
	sc := config.Proxy["http"]
	return http.ListenAndServe(sc.Addr, proxy)
}

func httpsServer(config *Config) error {
	proxy := newProxy(config.Auth.Addr)
	sc := config.Proxy["https"]
	certFile := sc.Params["cert"]
	keyFile := sc.Params["key"]
	return http.ListenAndServeTLS(sc.Addr, certFile, keyFile, proxy)
}
