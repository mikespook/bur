package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/mikespook/bur"
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
		log.Printf("%s\n\n", err)
		flag.Usage()
		return
	}
	bur.Serve(cfg)
}
