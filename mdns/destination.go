package mdns

import (
	"github.com/miekg/dns"
)

func (c *client) broadcastMessage(q *dns.Msg) error {
	buf, err := q.Pack()
	if err != nil {
		return err
	}
	if c.ipv4UnicastConn != nil {
		_, err = c.ipv4MulticastConn.Write(buf)
		if err != nil {
			return err
		}
	}
	if c.ipv6UnicastConn != nil {
		_, err = c.ipv6MulticastConn.Write(buf)
		if err != nil {
			return err
		}
	}
	return nil
}
