package mdns

import (
	"fmt"
	"strings"

	"github.com/miekg/dns"
)

func Scan(entries chan<- *dns.RR) error {
	c, err := newClient(true, false)
	if err != nil {
		return err
	}

	msgCh := make(chan *dns.Msg, 32)
	if c.use_ipv4 {
		// go c.recv(c.ipv4UnicastConn, msgCh)
		go c.recv(c.ipv4MulticastConn, msgCh)
	}
	if c.use_ipv6 {
		// go c.recv(c.ipv6UnicastConn, msgCh)
		go c.recv(c.ipv6MulticastConn, msgCh)
	}

	// inprogress := make(map[string]*ServiceEntry)

	for {
		select {
		case resp := <-msgCh:
			// var inp *ServiceEntry
			for _, answer := range append(resp.Answer, resp.Extra...) {
				entries <- &answer
				// fmt.Println("Receive answer \n", answer.String())
				// mdnsParser(answer, inp, inprogress)
				// if inp == nil {
				// 	continue
				// }
			}
		}
	}
}

func mdnsParser(raw dns.RR, inp *ServiceEntry, inprogress map[string]*ServiceEntry) {
	switch rr := raw.(type) {
	case *dns.PTR:
		// Create new entry for this
		inp = ensureName(inprogress, rr.Ptr)

	case *dns.SRV:
		// Check for a target mismatch
		if rr.Target != rr.Hdr.Name {
			alias(inprogress, rr.Hdr.Name, rr.Target)
		}

		// Get the port
		inp = ensureName(inprogress, rr.Hdr.Name)
		inp.Host = rr.Target
		inp.Port = int(rr.Port)

	case *dns.TXT:
		// Pull out the txt
		inp = ensureName(inprogress, rr.Hdr.Name)
		inp.Info = strings.Join(rr.Txt, "|")
		inp.InfoFields = rr.Txt
		inp.hasTXT = true

	case *dns.A:
		// Pull out the IP
		inp = ensureName(inprogress, rr.Hdr.Name)
		inp.Addr = rr.A // @Deprecated
		inp.AddrV4 = rr.A

	case *dns.AAAA:
		// Pull out the IP
		inp = ensureName(inprogress, rr.Hdr.Name)
		inp.Addr = rr.AAAA // @Deprecated
		inp.AddrV6 = rr.AAAA
	}
	fmt.Println("Receive answer \n", (*inp))
}
