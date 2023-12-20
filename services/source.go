package services

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"example.com/mdns-broadcast/mdns"
)

func SourceMain() {

	// Make a channel for results and start listening
	entriesCh := make(chan *mdns.ServiceEntry, 4)
	go func() {
		for entry := range entriesCh {
			fmt.Printf("Got new entry: %v\n", entry)
		}
	}()

	go server()

	// Start the lookup
	mdns.Lookup("_rfb._udp", entriesCh)
	close(entriesCh)
}

func server() {
	http.HandleFunc("/mdns", func(w http.ResponseWriter, r *http.Request) {
		mdnsEntries := mdns.GetMdnsEntries()
		// Marshal the mDNS entries to JSON
		jsonEntries, err := json.Marshal(mdnsEntries)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Set the content type and write the JSON to the response
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonEntries)
	})

	log.Fatal(http.ListenAndServe(":8081", nil))
}
