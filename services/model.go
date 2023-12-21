package services

import (
	"encoding/json"
	"fmt"

	"github.com/miekg/dns"
)

// MDNSRecord represents the JSON structure of a DNS resource record
type MDNSRecord struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Class string `json:"class"`
	TTL   uint32 `json:"ttl"`
	Data  string `json:"data"`
}

// SerializeMDNSRecords serializes a slice of dns.RR to JSON
func SerializeMDNSRecords(records []dns.RR) ([]byte, error) {
	var mdnsRecords []MDNSRecord
	for _, rr := range records {
		mdnsRecord := MDNSRecord{
			Name:  rr.Header().Name,
			Type:  dns.TypeToString[rr.Header().Rrtype],
			Class: dns.ClassToString[rr.Header().Class],
			TTL:   rr.Header().Ttl,
			Data:  rr.String(),
		}
		mdnsRecords = append(mdnsRecords, mdnsRecord)
	}
	jsonData, err := json.Marshal(mdnsRecords)
	return jsonData, err
}

// DeserializeMDNSRecords deserializes a JSON string to a slice of dns.RR
func DeserializeMDNSRecords(jsonData string) ([]dns.RR, error) {
	var mdnsRecords []MDNSRecord
	err := json.Unmarshal([]byte(jsonData), &mdnsRecords)
	if err != nil {
		return nil, err
	}

	var records []dns.RR
	for _, record := range mdnsRecords {
		rr, err := dns.NewRR(record.Data)
		if err != nil {
			fmt.Printf("Error parsing record: %v, input: %s", err, record.Data)
			continue
		}
		if rr == nil {
			fmt.Printf("Parsed record is nil, input: %s", record.Data)
			continue
		}
		records = append(records, rr)
	}
	return records, nil
}
