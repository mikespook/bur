package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/mikespook/bur"
	"github.com/mikespook/golib/log"
	"github.com/mikespook/golib/signal"
)

var config string

func init() {
	flag.Usage = func() {
		intro := `
Bur Proxy Server
`
		fmt.Println(intro)
		flag.PrintDefaults()
	}
	flag.StringVar(&config, "config", "", "Path to the configuration file")
	flag.Parse()
}

func main() {
	if config == "" {
		flag.Usage()
		return
	}
	cfg, err := bur.LoadConfig(config)
	if err != nil {
		log.Error(err)
		flag.Usage()
		return
	}
	if err := log.Init(cfg.Log.File, log.StrToLevel(cfg.Log.Level)); err != nil {
		log.Error(err)
	}
	go bur.Serve(cfg)
	sh := signal.NewHandler()
	sh.Bind(os.Interrupt, func() bool {
		bur.Close()
		return true
	})
	sh.Loop()
}
