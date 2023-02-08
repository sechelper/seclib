package network_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/miekg/dns"
	"github.com/sechelper/seclib/network"
	"log"
	"net"
)

func ExampleParseCIDR() {
	ipv4Addr, ipv4NetSegment, err := network.ParseCIDR("192.0.2.1/24")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(ipv4Addr)
	for i := range ipv4NetSegment {
		fmt.Println(ipv4NetSegment[i])
	}
}

func ExampleAddr_String() {
	addr := network.Addr{
		IP:   net.ParseIP("114.114.114.114"),
		Port: 53,
	}

	fmt.Println(addr)
}

func ExampleDns_Exchange() {
	var ips []net.IP = nil
	in, _, err := network.DefaultResolver.Exchange(dns.TypeA, "go-hacker-code.lab.secself.com")

	if err != nil || in.Answer == nil {
		log.Fatal(err)
	}

	ips = make([]net.IP, 0)
	for i := range in.Answer {
		if rr, ok := in.Answer[i].(*dns.A); ok {
			ips = append(ips, rr.A)
		}
	}
	fmt.Println(ips)
}

func ExampleDns_LookupIP() {
	ip, err := network.DefaultResolver.LookupIP("go-hacker-code.lab.secself.com")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(ip)
}

func ExampleDns_LookupCNAME() {
	cname, err := network.DefaultResolver.LookupCNAME("go-hacker-code.lab.secself.com")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(cname)
}

func ExampleRandomUserAgent() {
	for i := 0; i < 10; i++ {
		fmt.Println(network.RandomUserAgent())
	}
}

func ExampleIsDomainName() {
	fmt.Println(network.IsDomainName("secself.com"))        // true
	fmt.Println(network.IsDomainName("www.secself"))        // true
	fmt.Println(network.IsDomainName("secself.com/"))       // false
	fmt.Println(network.IsDomainName("http://secself.com")) // false
}

func ExampleHttpClient_Get() {
	c := network.HttpClient{}
	if reps, err := c.Get("https://go-hacker-code.lab.secself.com/"); err == nil {
		fmt.Println(reps.StatusCode)
	}
}

func ExampleHttpClient_Post() {
	c := network.HttpClient{}
	data := map[string]string{"username": "password"}
	body, _ := json.Marshal(data)
	if reps, err := c.Post("https://go-hacker-code.lab.secself.com/", "application/json", bytes.NewReader(body)); err == nil {
		fmt.Println(reps.StatusCode)
	}
}

func ExampleHttpClient_UploadFile() {
	c := network.HttpClient{}
	file1 := network.UPFile{Path: "file1.txt",
		Field: "file"}

	file2 := network.UPFile{Path: "file2.txt",
		Field: "file1"}

	file3 := network.UPFile{Path: "file3.txt",
		Field: "file2"}

	_, err := c.UploadFile("http://127.0.0.1:8080/upladMore", []network.UPFile{file1, file2, file3})
	if err != nil {
		fmt.Println(err)
		return
	}
}
