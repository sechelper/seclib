package network

import (
	"github.com/miekg/dns"
	"net"
	"strings"
	"time"
)

type Msg dns.Msg

type Dns struct {
	Ns      string
	Msg     *Msg
	Timeout time.Duration
}

func NewDefaultMsg() *Msg {
	msg := new(Msg)
	msg.Id = dns.Id()
	msg.RecursionDesired = true
	msg.Truncated = true
	return msg
}

var DefaultResolver = &Dns{
	Msg:     NewDefaultMsg(),
	Ns:      "223.6.6.6:53",
	Timeout: 3 * time.Second,
}

func (_dns *Dns) Exchange(dnsType uint16, domain string) (r *dns.Msg, rtt time.Duration, err error) {
	if !strings.HasSuffix(domain, ".") {
		domain += "."
	}
	_dns.Msg.Question = make([]dns.Question, 1)
	_dns.Msg.Question[0] = dns.Question{Name: domain,
		Qtype: dnsType, Qclass: dns.ClassINET}
	c := new(dns.Client)
	c.Timeout = _dns.Timeout
	return c.Exchange((*dns.Msg)(_dns.Msg), _dns.Ns)
}

// LookupIP looks up host using the local resolver.
// It returns a slice of that host's IPv4 and IPv6 addresses.
func (_dns *Dns) LookupIP(host string) ([]net.IP, error) {
	var ips []net.IP = nil
	in, _, err := DefaultResolver.Exchange(dns.TypeA, host)

	if err != nil || in.Answer == nil {
		return nil, err
	}

	ips = make([]net.IP, 0)
	for i := range in.Answer {
		if rr, ok := in.Answer[i].(*dns.A); ok {
			ips = append(ips, rr.A)
		}
	}
	return ips, nil
}

func (_dns *Dns) LookupCNAME(host string) (cname string, err error) {
	in, _, err := DefaultResolver.Exchange(dns.TypeCNAME, host)

	if err != nil || in.Answer == nil {
		return "", err
	}
	if rr, ok := in.Answer[0].(*dns.CNAME); ok {
		return rr.Target, nil
	}
	return "", nil
}
