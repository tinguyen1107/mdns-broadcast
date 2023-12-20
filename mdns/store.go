package mdns

import (
	"sync"
)

// Global variable
var (
	mdnsEntries []string
	mutex       sync.Mutex
)

// Writes data to the global variable
func AddMdnsEntries(entry string) {
	mutex.Lock() // Lock the mutex before writing
	mdnsEntries = append(mdnsEntries, entry)
	mutex.Unlock() // Unlock the mutex after writing
}

// Reads data from the global variable
func GetMdnsEntries() []string {
	mutex.Lock() // Lock the mutex before reading
	entriesCopy := make([]string, len(mdnsEntries))
	copy(entriesCopy, mdnsEntries)
	mutex.Unlock() // Unlock the mutex after reading
	return entriesCopy
}
