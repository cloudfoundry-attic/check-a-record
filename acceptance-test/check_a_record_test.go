package acceptance_test

import (
	"net"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("check-a-record", func() {
	AfterEach(func() {
		dnsServer.DeregisterAllRecords()
	})

	Context("when A records exist", func() {
		It("exits 0 and prints only the A records", func() {
			dnsServer.RegisterARecord("domain-with-a-and-mx", net.IP{1, 2, 3, 4})
			dnsServer.RegisterMXRecord("domain-with-a-and-mx", "some-mail-server.", 0)

			session := checkARecord([]string{"domain-with-a-and-mx"})
			Eventually(session, time.Minute).Should(gexec.Exit(0))

			Expect(session.Out.Contents()).To(ContainSubstring("1.2.3.4"))
			Expect(session.Out.Contents()).NotTo(ContainSubstring("some-mail-server."))
		})
	})

	Context("when no A records exist", func() {
		It("exits 1 and prints an error", func() {
			dnsServer.RegisterMXRecord("domain-with-two-mx-records", "some-mail-server.", 0)
			dnsServer.RegisterMXRecord("domain-with-two-mx-records", "another-mail-server.", 1)

			session := checkARecord([]string{"domain-with-two-mx-records"})
			Eventually(session, time.Minute).Should(gexec.Exit(1))

			Expect(session.Err.Contents()).To(ContainSubstring("lookup domain-with-two-mx-records on 127.0.0.1:53: no such host"))
		})
	})

	Context("when the domain does not exist at all", func() {
		It("exits 1 and prints an error", func() {
			session := checkARecord([]string{"nonexistent-domain"})
			Eventually(session, time.Minute).Should(gexec.Exit(1))

			Expect(session.Err.Contents()).To(ContainSubstring("lookup nonexistent-domain on 127.0.0.1:53: no such host"))
		})
	})
})
