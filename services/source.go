package services

import (
	"fmt"
	"time"

	"github.com/miekg/dns"

	"example.com/mdns-broadcast/mdns"
)

func SourceMain() {
	// Make a channel for results and start listening
	entriesCh := make(chan *dns.Msg, 4)
	go func() {
		for entry := range entriesCh {
			fmt.Println("Got new entry: ", *entry)
			AddMdnsEntries(*entry)
		}
	}()
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		for range ticker.C {
			UploadMdnsEntries()
		}
	}()

	// Start the scan
	mdns.Scan(entriesCh)
	close(entriesCh)
}
