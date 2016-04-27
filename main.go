package main

import (
	"fmt"
	"os"

	"github.com/cloudfoundry-incubator/check-a-record/dns"
	miekgdns "github.com/miekg/dns"
)

func main() {
	domain := os.Args[1]

	dnsResolver := dns.NewDNSResolver(&miekgdns.Client{}, os.Args[2])
	ips, err := dnsResolver.ResolveARecord(domain)
	if err != nil {
		fail(err.Error())
	}

	if len(ips) == 0 {
		fail(fmt.Sprintf("No A records exist for %s.", domain))
	}

	fmt.Printf("%v", ips)
}

func fail(message string) {
	fmt.Fprintf(os.Stderr, "%s\n", message)
	os.Exit(1)
}
