package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/miekg/dns"
)

// Global variable
var (
	mdnsEntries []dns.RR
	mutex       sync.Mutex
)

// Writes data to the global variable
func AddMdnsEntries(entry dns.RR) {
	mutex.Lock() // Lock the mutex before writing
	mdnsEntries = append(mdnsEntries, entry)
	mutex.Unlock() // Unlock the mutex after writing
}

// Reads data from the global variable
func GetMdnsEntries() []dns.RR {
	mutex.Lock() // Lock the mutex before reading
	entriesCopy := make([]dns.RR, len(mdnsEntries))
	copy(entriesCopy, mdnsEntries)
	mutex.Unlock() // Unlock the mutex after reading
	return entriesCopy
}

// Uploads mDNS entries to a server
func UploadMdnsEntries() {
	fmt.Println("MDNS UPLOAD")
	entries := GetMdnsEntries()
	jsonData, err := json.Marshal(entries)
	if err != nil {
		// Handle error
		return
	}

	// Replace with the actual server URL
	url := "http://localhost:8081/mdns-entries"
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		// Handle error
		return
	}
	defer resp.Body.Close()

	// Optionally, handle the response from the server
}
