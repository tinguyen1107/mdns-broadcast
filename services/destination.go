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
		data, err := DeserializeMDNSMessageList(string(body))
		if err != nil {
			http.Error(w, "Error parsing request body", http.StatusInternalServerError)
			return
		}

		broadcastMDNS(data)
		// fmt.Println("Request Body:", data)
	})

	log.Fatal(http.ListenAndServe(":8081", nil))
}

func broadcastMDNS(msgs []*dns.Msg) {
	// mDNS uses the 224.0.0.251 multicast address and port 5353
	const mdnsAddress = "224.0.0.251:5353"

	for _, m := range msgs {
		// Prepare the message for transmission
		msgData, err := m.Pack()
		if err != nil {
			log.Printf("Failed to pack mDNS message: %v", err)
			continue // move to the next message
		}

		// Set up UDP connection for mDNS
		conn, err := net.Dial("udp", mdnsAddress)
		if err != nil {
			log.Printf("Failed to set up UDP connection for mDNS: %v", err)
			continue // move to the next message
		}

		// Send the mDNS response message
		_, err = conn.Write(msgData)
		if err != nil {
			log.Printf("Failed to write mDNS message: %v", err)
		} else {
			// Log success if no error
			log.Printf("Successfully broadcasted mDNS message")
		}

		// Close the connection
		conn.Close()
	}
}
