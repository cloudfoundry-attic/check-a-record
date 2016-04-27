package acceptance_test

import (
	"net"
	"time"

	"github.com/cloudfoundry-incubator/check-a-record/acceptance-test/dnsserver"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("check-a-record", func() {
	var server dnsserver.Server

	BeforeEach(func() {
		server = dnsserver.NewServer()
		server.Start()
	})

	AfterEach(func() {
		server.Stop()
	})

	Context("when A records exist", func() {
		It("exits 0 and prints only the A records", func() {
			server.RegisterCNAMERecord("some-domain", "some-cname.")
			server.RegisterARecord("some-domain", net.IP{1, 2, 3, 4})

			session := checkARecord([]string{"some-domain", server.URL()})
			Eventually(session, time.Minute).Should(gexec.Exit(0))

			Expect(session.Out.Contents()).To(ContainSubstring("1.2.3.4"))
			Expect(session.Out.Contents()).NotTo(ContainSubstring("some-cname."))
		})
	})

	Context("when no A records exist", func() {
		It("exits 1 and prints an error", func() {
			server.RegisterCNAMERecord("some-domain", "some-cname.")
			server.RegisterCNAMERecord("some-domain", "another-cname.")

			session := checkARecord([]string{"some-domain", server.URL()})
			Eventually(session, time.Minute).Should(gexec.Exit(1))

			Expect(session.Err.Contents()).To(ContainSubstring("No A records exist for some-domain"))
		})
	})

	Context("failure cases", func() {
		It("exits 1 and prints an error when dns server doesn't exist", func() {
			session := checkARecord([]string{"some-domain", "127.0.0.1:9999"})
			Eventually(session, time.Minute).Should(gexec.Exit(1))

			Expect(session.Err.Contents()).To(ContainSubstring("read: connection refused"))
		})
	})
})
