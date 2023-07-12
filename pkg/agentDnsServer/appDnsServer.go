// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package agentDnsServer

import (
	"crypto/tls"
	"github.com/miekg/dns"
	"go.uber.org/zap"
	"net"
	"sync/atomic"
)

var RequestDnsCounts int64 = 0

func SetupAppDnsServer(rootLogger *zap.Logger, tlsCert, tlsKey string) {
	logger := rootLogger.Named("app dns server")
	logger.Sugar().Infof("setup app dns Server")
	cert, err := tls.LoadX509KeyPair(tlsCert, tlsKey)
	if err != nil {
		logger.Sugar().Fatalf("failed load cert key from path %s %s", tlsCert, tlsKey)
	}
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	handler := dns.HandlerFunc(func(w dns.ResponseWriter, r *dns.Msg) {
		m := new(dns.Msg)
		m.SetReply(r)
		m.Answer = []dns.RR{
			&dns.A{
				Hdr: dns.RR_Header{Name: r.Question[0].Name, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 3600},
				A:   net.ParseIP("127.0.0.1"),
			},
			&dns.AAAA{
				Hdr:  dns.RR_Header{Name: r.Question[0].Name, Rrtype: dns.TypeAAAA, Class: dns.ClassINET, Ttl: 3600},
				AAAA: net.ParseIP("::1"),
			},
		}
		atomic.AddInt64(&RequestDnsCounts, 1)
		w.WriteMsg(m)
	})

	tlsServer := &dns.Server{
		Addr:      ":853",
		Net:       "tcp-tls",
		TLSConfig: tlsConfig,
		Handler:   handler,
	}
	udpServer := &dns.Server{
		Addr:    ":53",
		Net:     "udp",
		Handler: handler,
	}

	tcpServer := &dns.Server{
		Addr:    ":53",
		Net:     "tcp",
		Handler: handler,
	}
	go func() {
		if err := tcpServer.ListenAndServe(); err != nil {
			logger.Sugar().Fatalf("dns tcp server , err: %v ", err)
		}
	}()

	go func() {
		if err := udpServer.ListenAndServe(); err != nil {
			logger.Sugar().Fatalf("dns udp server , err: %v ", err)
		}
	}()

	go func() {
		if err := tlsServer.ListenAndServe(); err != nil {
			logger.Sugar().Fatalf("dns tcp-tls server , err: %v ", err)
		}
	}()
	logger.Sugar().Infof("dns udp server listien %s", udpServer.Addr)
	logger.Sugar().Infof("dns tcp server listien %s", tcpServer.Addr)
	logger.Sugar().Infof("dns tcp-tls server listien %s", tcpServer.Addr)
}
