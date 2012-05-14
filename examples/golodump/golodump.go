package main

import (
	"github.com/surma/golo"
	"net"
	"flag"
	"log"
)

var (
	addr = flag.String("a", ":7770", "Address to listen on")
)

func main() {
	flag.Parse()
	addr, e := net.ResolveUDPAddr("udp4", *addr)
	if e != nil {
		log.Fatalf("Could not resolve addr \"%s\": %s", *addr, e)
	}
	c, e := net.ListenUDP("udp4", addr)
	if e != nil {
		log.Fatalf("Could not listen: %s", e)
	}

	// Maxiumum UPD package payload length
	// due to IPv4 and stuff
	b := make([]byte, 65507)
	for {
		n, _, e := c.ReadFrom(b)
		if e != nil {
			log.Printf("Read failed: %s", e)
			continue
		}
		msg, e := golo.Deserialize(b[0:n])
		if e != nil {
			log.Printf("Received invalid package: %s", e)
		}
		log.Printf("Received: %s", msg)
	}
}
