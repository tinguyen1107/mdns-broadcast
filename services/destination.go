package services

import (
	"fmt"
	"os"

	"example.com/mdns-broadcast/mdns"
)

func DestinationMain() {
	// Setup our service export
	host, _ := os.Hostname()
	info := []string{"My awesome service"}
	service, err := mdns.NewMDNSService(host, "_foobar._tcp", "", "", 8000, nil, info)

	if err != nil {
		fmt.Println("Create new mdns get error", err)
	}
	fmt.Println("Create new mdns pass")

	// Create the mDNS server, defer shutdown
	server, err := mdns.NewServer(&mdns.Config{Zone: service})
	if err != nil {
		fmt.Println("Create new Server get error", err)
	}
	fmt.Println("Create new Server pass")
	defer server.Shutdown()
}
