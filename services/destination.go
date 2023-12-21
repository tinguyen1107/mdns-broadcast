package services

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"example.com/mdns-broadcast/mdns"
)

// DestinationMain
// using for receive the grpc request from source and broadcast mdns to local network
func DestinationMain() {
	// Setup our service export
	host, _ := os.Hostname()
	info := []string{"My awesome service"}
	service, err := mdns.NewMDNSService(host, "_foobar._tcp", "", "", 8000, nil, info)

	if err != nil {
		fmt.Println("Create new mdns get error", err)
	}
	fmt.Println("Create new mdns pass")

	// server using for receiving mdns from source
	server(service)

	// Create the mDNS server, defer shutdown
	server, err := mdns.NewServer(&mdns.Config{Zone: service})
	if err != nil {
		fmt.Println("Create new Server get error", err)
	}
	fmt.Println("Create new Server pass")
	defer server.Shutdown()
}

func server(service *mdns.MDNSService) {
	http.HandleFunc("/mdns-entries", func(w http.ResponseWriter, r *http.Request) {
		// mdnsEntries := mdns.GetMdnsEntries()
		// // Marshal the mDNS entries to JSON
		// jsonEntries, err := json.Marshal(mdnsEntries)
		// if err != nil {
		// 	http.Error(w, err.Error(), http.StatusInternalServerError)
		// 	return
		// }
		//
		// // Set the content type and write the JSON to the response
		// w.Header().Set("Content-Type", "application/json")
		// w.Write(jsonEntries)
		// Read the request body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		// Convert the body to a string and print it
		var data any
		json.Unmarshal(body, &data)
		fmt.Println("Request Body:", data)
	})

	log.Fatal(http.ListenAndServe(":8081", nil))
}
