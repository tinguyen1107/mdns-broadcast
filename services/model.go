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

// MDNSMessage represents the JSON structure of a DNS message
type MDNSMessage struct {
	MsgHdr    dns.MsgHdr     `json:"header"`
	Questions []dns.Question `json:"questions"`
	Answers   []MDNSRecord   `json:"answers"`
	Ns        []MDNSRecord   `json:"authorities"`
	Extra     []MDNSRecord   `json:"additional"`
}

// SerializeMDNSMessageList serializes a slice of dns.Msg to JSON
func SerializeMDNSMessageList(msgs []dns.Msg) ([]byte, error) {
	var mdnsMsgs []MDNSMessage
	for _, msg := range msgs {
		mdnsMsg := MDNSMessage{
			MsgHdr:    msg.MsgHdr,
			Questions: msg.Question,
			Answers:   convertRRSliceToMDNSRecord(msg.Answer),
			Ns:        convertRRSliceToMDNSRecord(msg.Ns),
			Extra:     convertRRSliceToMDNSRecord(msg.Extra),
		}
		mdnsMsgs = append(mdnsMsgs, mdnsMsg)
	}

	jsonData, err := json.Marshal(mdnsMsgs)
	return jsonData, err
}

// SerializeMDNSMessage serializes a dns.Msg to JSON
func SerializeMDNSMessage(msg *dns.Msg) ([]byte, error) {
	mdnsMsg := MDNSMessage{
		MsgHdr:    msg.MsgHdr,
		Questions: msg.Question,
		Answers:   convertRRSliceToMDNSRecord(msg.Answer),
		Ns:        convertRRSliceToMDNSRecord(msg.Ns),
		Extra:     convertRRSliceToMDNSRecord(msg.Extra),
	}

	jsonData, err := json.Marshal(mdnsMsg)
	return jsonData, err
}

func convertRRSliceToMDNSRecord(rrs []dns.RR) []MDNSRecord {
	var mdnsRecords []MDNSRecord
	for _, rr := range rrs {
		mdnsRecord := MDNSRecord{
			Name:  rr.Header().Name,
			Type:  dns.TypeToString[rr.Header().Rrtype],
			Class: dns.ClassToString[rr.Header().Class],
			TTL:   rr.Header().Ttl,
			Data:  rr.String(),
		}
		mdnsRecords = append(mdnsRecords, mdnsRecord)
	}
	return mdnsRecords
}

// DeserializeMDNSMessageList deserializes a JSON string to a slice of dns.Msg
func DeserializeMDNSMessageList(jsonData string) ([]*dns.Msg, error) {
	var mdnsMsgs []MDNSMessage
	err := json.Unmarshal([]byte(jsonData), &mdnsMsgs)
	if err != nil {
		return nil, err
	}

	var msgs []*dns.Msg
	for _, mdnsMsg := range mdnsMsgs {
		msg := new(dns.Msg)
		msg.MsgHdr = mdnsMsg.MsgHdr
		msg.Question = mdnsMsg.Questions
		msg.Answer = convertMDNSRecordToRRSlice(mdnsMsg.Answers)
		msg.Ns = convertMDNSRecordToRRSlice(mdnsMsg.Ns)
		msg.Extra = convertMDNSRecordToRRSlice(mdnsMsg.Extra)

		msgs = append(msgs, msg)
	}
	return msgs, nil
}

// DeserializeMDNSMessage deserializes a JSON string to a dns.Msg
func DeserializeMDNSMessage(jsonData string) (*dns.Msg, error) {
	var mdnsMsg MDNSMessage
	err := json.Unmarshal([]byte(jsonData), &mdnsMsg)
	if err != nil {
		return nil, err
	}

	msg := new(dns.Msg)
	msg.MsgHdr = mdnsMsg.MsgHdr
	msg.Question = mdnsMsg.Questions
	msg.Answer = convertMDNSRecordToRRSlice(mdnsMsg.Answers)
	msg.Ns = convertMDNSRecordToRRSlice(mdnsMsg.Ns)
	msg.Extra = convertMDNSRecordToRRSlice(mdnsMsg.Extra)

	return msg, nil
}

func convertMDNSRecordToRRSlice(mdnsRecords []MDNSRecord) []dns.RR {
	var rrs []dns.RR
	for _, record := range mdnsRecords {
		rr, err := dns.NewRR(record.Data)
		if err != nil {
			fmt.Printf("Error parsing record: %v, input: %s\n", err, record.Data)
			continue
		}
		if rr == nil {
			fmt.Printf("Parsed record is nil, input: %s\n", record.Data)
			continue
		}
		rrs = append(rrs, rr)
	}
	return rrs
}
