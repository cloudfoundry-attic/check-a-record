package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	domain := os.Args[1]
	ips, err := net.LookupIP(domain)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}

	filteredIPV4s := ipV4s(ips)
	fmt.Printf("%+v", filteredIPV4s)
}

func ipV4s(ips []net.IP) []net.IP {
	ipV4s := []net.IP{}
	for _, ip := range ips {
		if ipV4 := ip.To4(); ipV4 != nil {
			ipV4s = append(ipV4s, ipV4)
		}
	}

	return ipV4s
}
