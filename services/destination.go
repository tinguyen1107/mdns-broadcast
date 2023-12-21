package services

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/miekg/dns"

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
	// c, _ := mdns.NewClient(true, false)
	http.HandleFunc("/mdns-entries", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		// Convert the body to a string and print it
		data, err := DeserializeMDNSRecords(string(body))
		if err != nil {
			http.Error(w, "Error parsing request body", http.StatusInternalServerError)
			return
		}

		broadcastMDNS(data)
		// fmt.Println("Request Body:", data)
	})

	log.Fatal(http.ListenAndServe(":8081", nil))
}

func broadcastMDNS(records []dns.RR) {
	m := new(dns.Msg)
	m.Response = true  // This is a response message
	m.Answer = records // Add the parsed records to the answer section

	msgData, err := m.Pack()
	if err != nil {
		log.Printf("Failed to pack mDNS message: %v", err)
		return
	}

	// mDNS uses the 224.0.0.251 multicast address and port 5353
	const mdnsAddress = "224.0.0.251:5353"

	conn, err := net.Dial("udp", mdnsAddress)
	if err != nil {
		log.Printf("Failed to set up UDP connection for mDNS: %v", err)
		return
	}
	defer conn.Close()

	// Send the mDNS response message
	// err = m.WriteMsg(dns.NewMsgConn(conn))

	log.Printf(" mDNS message: %v", msgData)
	conn.Write(msgData)
	if err != nil {
		log.Printf("Failed to write mDNS message: %v", err)
		return
	}
}

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
