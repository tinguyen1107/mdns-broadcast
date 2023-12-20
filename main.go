package main

import (
	"example.com/mdns-broadcast/mdns"
	"fmt"
)

func main() {
	// Make a channel for results and start listening
	entriesCh := make(chan *mdns.ServiceEntry, 4)
	go func() {
		for entry := range entriesCh {
			fmt.Printf("Got new entry: %v\n", entry)
		}
	}()

	// Start the lookup
	mdns.Lookup("_rfb._udp", entriesCh)
	close(entriesCh)
}
