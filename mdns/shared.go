package mdns

import "net"

const (
	ipv4mdns              = "224.0.0.251"
	ipv6mdns              = "ff02::fb"
	mdnsPort              = 5353
	forceUnicastResponses = false
)

var (
	ipv4Addr = &net.UDPAddr{
		IP:   net.ParseIP(ipv4mdns),
		Port: mdnsPort,
	}
	ipv6Addr = &net.UDPAddr{
		IP:   net.ParseIP(ipv6mdns),
		Port: mdnsPort,
	}
)
