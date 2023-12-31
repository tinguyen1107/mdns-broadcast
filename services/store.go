package services

import (
	"bytes"
	"fmt"
	"net/http"
	"sync"

	"github.com/miekg/dns"
)

// Global variable
var (
	mdnsEntries []dns.Msg
	mutex       sync.Mutex
)

// Writes data to the global variable
func AddMdnsEntries(entry dns.Msg) {
	mutex.Lock() // Lock the mutex before writing
	mdnsEntries = append(mdnsEntries, entry)
	mutex.Unlock() // Unlock the mutex after writing
}

// Reads data from the global variable
func GetMdnsEntries() []dns.Msg {
	mutex.Lock() // Lock the mutex before reading
	entriesCopy := make([]dns.Msg, len(mdnsEntries))
	copy(entriesCopy, mdnsEntries)
	mutex.Unlock() // Unlock the mutex after reading
	return entriesCopy
}

func GetAndClearMdnsEntries() []dns.Msg {
	mutex.Lock() // Lock the mutex before reading
	entriesCopy := make([]dns.Msg, len(mdnsEntries))
	copy(entriesCopy, mdnsEntries)
	mdnsEntries = []dns.Msg{}
	mutex.Unlock() // Unlock the mutex after reading
	return entriesCopy
}

// Uploads mDNS entries to a server
func UploadMdnsEntries() {
	fmt.Println("MDNS UPLOAD")
	entries := GetAndClearMdnsEntries()
	jsonData, err := SerializeMDNSMessageList(entries)
	if err != nil {
		// Handle error
		return
	}

	// Replace with the actual server URL
	url := "http://192.168.128.132:8081/mdns-entries"
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		// Handle error
		return
	}
	defer resp.Body.Close()

	// Optionally, handle the response from the server
}
