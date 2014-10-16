package main

import (
	"flag"
	"net/rpc"
	"net/url"

	"log"

	"github.com/mikespook/bur"
)

var (
	server, username string
)

func init() {
	flag.StringVar(&server, "server", "", "URI to the Bur server")
	flag.StringVar(&username, "user", "", "Name to look for")
	flag.Parse()
}

func main() {
	if server == "" {
		flag.Usage()
		return
	}
	u, err := url.Parse(server)
	if err != nil {
		log.Println(err)
		return
	}
	client, err := rpc.Dial(u.Scheme, u.Host)
	if err != nil {
		log.Println(err)
		return
	}
	if username != "" {
		var user bur.User
		if err := client.Call("Bur.User", username, &user); err != nil {
			log.Println(err)
			return
		}
		log.Printf("%v\n", user)
		return
	}
	var state bur.State
	if err := client.Call("Bur.State", "", &state); err != nil {
		log.Println(err)
		return
	}
	log.Printf("%v\n", state)
}
