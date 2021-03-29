package main

import (
	"flag"
	"log"
)

func main() {
	var pemPath string
	var keyPath string
	var proto string

	flag.StringVar(&pemPath, "pem", "server.pem", "path to pem file")
	flag.StringVar(&keyPath, "key", "server.key", "path to key file")
	flag.StringVar(&proto, "proto", "http", "Proxy protocol (http or https)")
	flag.Parse()

	if proto != "http" && proto != "https" {
		log.Fatal("Protocol must be either http or https")
	}
	log.Println(pemPath)
	log.Println(keyPath)
	log.Println(proto)
	StartProxy(pemPath, keyPath, proto)

}
