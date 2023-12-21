package mdns

import (
	"fmt"

	"github.com/miekg/dns"
)

type Client struct {
	c client
}

func NewClient(v4 bool, v6 bool) (*Client, error) {
	c, err := newClient(true, false)
	if err != nil {
		return nil, fmt.Errorf("failed to bind to any multicast udp port")
	}

	client := &Client{
		c: *c,
	}
	return client, nil
}

func (c *Client) BroadcastMessage(q *dns.Msg) error {
	buf, err := q.Pack()
	if err != nil {
		return err
	}
	if c.c.ipv4UnicastConn != nil {
		_, err = c.c.ipv4MulticastConn.Write(buf)
		if err != nil {
			return err
		}
	}
	if c.c.ipv6UnicastConn != nil {
		_, err = c.c.ipv6MulticastConn.Write(buf)
		if err != nil {
			return err
		}
	}
	return nil
}
