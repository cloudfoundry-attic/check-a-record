package fakes

import (
	"time"

	miekgdns "github.com/miekg/dns"
)

type DNSClient struct {
	ExchangeCall struct {
		Returns struct {
			Message *miekgdns.Msg
			Error   error
		}
		Receives struct {
			Message   *miekgdns.Msg
			DNSServer string
		}
	}
}

func (d *DNSClient) Exchange(message *miekgdns.Msg, dnsServer string) (*miekgdns.Msg, time.Duration, error) {
	d.ExchangeCall.Receives.Message = message
	d.ExchangeCall.Receives.DNSServer = dnsServer
	return d.ExchangeCall.Returns.Message, time.Second, d.ExchangeCall.Returns.Error
}
