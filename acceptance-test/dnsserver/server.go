package dnsserver

import (
	"fmt"
	"net"

	miekgdns "github.com/miekg/dns"
)

type Server struct {
	server  miekgdns.Server
	records map[string][]miekgdns.RR
}

func NewServer() Server {
	server := Server{
		records: make(map[string][]miekgdns.RR),
	}

	server.server = miekgdns.Server{
		Addr:    ":8053",
		Net:     "udp",
		Handler: miekgdns.HandlerFunc(server.handleDNSRequest),
	}

	return server
}

func (s *Server) Start() error {
	go func() {
		err := s.server.ListenAndServe()

		if err != nil {
			fmt.Println(err.Error())
		}
	}()

	return nil
}

func (s *Server) Stop() {
	s.server.Shutdown()
}

func (Server) URL() string {
	return "127.0.0.1:8053"
}

func (s Server) RegisterARecord(domainName string, ipAddress net.IP) {
	s.registerRecord(domainName, &miekgdns.A{
		Hdr: s.header(domainName, miekgdns.TypeA),
		A:   ipAddress,
	})
}

func (s Server) RegisterCNAMERecord(domainName string, target string) {
	s.registerRecord(domainName, &miekgdns.CNAME{
		Hdr:    s.header(domainName, miekgdns.TypeCNAME),
		Target: target,
	})
}

func (s Server) handleDNSRequest(responseWriter miekgdns.ResponseWriter, requestMessage *miekgdns.Msg) {
	responseMessage := new(miekgdns.Msg)
	responseMessage.SetReply(requestMessage)

	resourceRecords, recordExists := s.records[requestMessage.Question[0].Name]

	if recordExists {
		responseMessage.Answer = make([]miekgdns.RR, len(resourceRecords))
		for i, resourceRecord := range resourceRecords {
			responseMessage.Answer[i] = resourceRecord
		}
	}

	responseWriter.WriteMsg(responseMessage)
}

func (s Server) registerRecord(domainName string, resourceRecord miekgdns.RR) {
	_, exists := s.records[domainName+"."]

	if !exists {
		s.records[domainName+"."] = []miekgdns.RR{}
	}

	s.records[domainName+"."] = append(s.records[domainName+"."], resourceRecord)
}

func (s Server) header(domainName string, resourceRecordType uint16) miekgdns.RR_Header {
	return miekgdns.RR_Header{
		Name:   domainName + ".",
		Rrtype: resourceRecordType,
		Class:  miekgdns.ClassINET,
		Ttl:    0,
	}
}
