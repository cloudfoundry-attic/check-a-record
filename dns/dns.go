package dns

import (
	"time"

	miekgdns "github.com/miekg/dns"
)

type dnsClient interface {
	Exchange(*miekgdns.Msg, string) (*miekgdns.Msg, time.Duration, error)
}

type DNSResolver struct {
	dnsClient        dnsClient
	dnsServerAddress string
}

func NewDNSResolver(dnsClient dnsClient, dnsServerAddress string) DNSResolver {
	return DNSResolver{
		dnsClient:        dnsClient,
		dnsServerAddress: dnsServerAddress,
	}
}

func (d DNSResolver) ResolveARecord(hostname string) ([]string, error) {
	requestMessage := miekgdns.Msg{}
	requestMessage.SetQuestion(hostname+".", miekgdns.TypeA)

	responseMessage, _, err := d.dnsClient.Exchange(&requestMessage, d.dnsServerAddress)
	if err != nil {
		return []string{}, err
	}

	var ips []string
	for _, answer := range responseMessage.Answer {
		ARecord, isARecord := answer.(*miekgdns.A)
		if isARecord {
			ips = append(ips, ARecord.A.String())
		}
	}

	return ips, nil
}
