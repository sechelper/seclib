package network

import (
	"encoding/binary"
	"fmt"
	"net"
)

type Addr struct {
	IP   net.IP
	Port int
}

func (addr Addr) String() string {
	return fmt.Sprintf("%s:%d", addr.IP, addr.Port)
}

var PRIVATE_IPV4 = [3]string{
	"10.0.0.0/8",
	"172.16.0.0/12",
	"192.168.0.0/16",
}

// ParseCIDR parses s as a CIDR notation IP address and prefix length,
// like "192.0.2.0/24" or "2001:db8::/32", as defined in
// RFC 4632 and RFC 4291.
//
// It returns the IP address and the network implied by the IP and
// prefix length.
// For example, ParseCIDR("192.0.2.1/24") returns the IP address
// 192.0.2.1 and the network segment.
// TODO support ipv6
func ParseCIDR(s string) (net.IP, []net.IP, error) {
	ips := make([]net.IP, 0)
	ipv4Addr, ipv4Net, err := net.ParseCIDR(s)
	if err != nil {
		return nil, nil, err
	}

	// convert IPNet struct mask and address to uint32
	// network is BigEndian
	mask := binary.BigEndian.Uint32(ipv4Net.Mask)
	start := binary.BigEndian.Uint32(ipv4Net.IP)

	// find the final address
	finish := (start & mask) | (mask ^ 0xffffffff)

	// loop through addresses as uint32
	for i := start; i <= finish; i++ {
		// convert back to net.IP
		ip := make(net.IP, 4)
		binary.BigEndian.PutUint32(ip, i)
		ips = append(ips, ip)
	}
	return ipv4Addr, ips, nil
}
