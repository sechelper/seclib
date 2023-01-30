package network_test

import (
	"fmt"
	"github.com/sechelper/seclib/network"
	"log"
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
