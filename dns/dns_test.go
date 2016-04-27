package dns_test

import (
	"errors"
	"net"

	"github.com/cloudfoundry-incubator/check-a-record/dns"
	"github.com/cloudfoundry-incubator/check-a-record/fakes"
	miekgdns "github.com/miekg/dns"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("DNS", func() {
	Describe("ResolveARecord", func() {
		var (
			dnsResolver   dns.DNSResolver
			fakeDNSClient *fakes.DNSClient
		)

		BeforeEach(func() {
			fakeDNSClient = &fakes.DNSClient{}
			dnsResolver = dns.NewDNSResolver(fakeDNSClient, "some-dns-server")
		})

		It("returns a list IPv4 for hostname", func() {
			fakeDNSClient.ExchangeCall.Returns.Message = &miekgdns.Msg{
				Answer: []miekgdns.RR{
					&miekgdns.A{
						A: net.IP{8, 8, 8, 8},
					},
				},
			}

			ips, err := dnsResolver.ResolveARecord("google.com")
			Expect(err).NotTo(HaveOccurred())

			Expect(fakeDNSClient.ExchangeCall.Receives.Message.Question[0]).To(Equal(miekgdns.Question{
				Name:   "google.com.",
				Qtype:  1,
				Qclass: 1,
			}))
			Expect(fakeDNSClient.ExchangeCall.Receives.DNSServer).To(Equal("some-dns-server"))

			Expect(ips).To(ConsistOf("8.8.8.8"))
		})

		It("ignores non-A records returned by the dns client", func() {
			fakeDNSClient.ExchangeCall.Returns.Message = &miekgdns.Msg{
				Answer: []miekgdns.RR{
					&miekgdns.CNAME{
						Target: "some-cname",
					},
					&miekgdns.A{
						A: net.IP{8, 8, 8, 8},
					},
				},
			}

			ips, err := dnsResolver.ResolveARecord("google.com")
			Expect(err).NotTo(HaveOccurred())

			Expect(ips).To(ConsistOf("8.8.8.8"))
		})

		Describe("failure cases", func() {
			It("returns an error if dns client fails", func() {
				fakeDNSClient.ExchangeCall.Returns.Error = errors.New("something bad happened")

				_, err := dnsResolver.ResolveARecord("google.com")
				Expect(err).To(MatchError("something bad happened"))
			})
		})
	})
})
