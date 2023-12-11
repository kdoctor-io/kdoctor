// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package agentDnsServer

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"strings"

	"github.com/miekg/dns"
	"go.uber.org/zap"

	"github.com/kdoctor-io/kdoctor/pkg/agentHttpServer"
	k8sObjManager "github.com/kdoctor-io/kdoctor/pkg/k8ObjManager"
	"github.com/kdoctor-io/kdoctor/pkg/types"
)

func SetupAppDnsServer(rootLogger *zap.Logger, tlsCert, tlsKey string) {
	logger := rootLogger.Named("app dns server")
	logger.Sugar().Infof("setup app dns Server")
	var tlsConfig *tls.Config
	var err error
	if tlsCert != "" && tlsKey != "" {
		cert, err := tls.LoadX509KeyPair(tlsCert, tlsKey)
		if err != nil {
			logger.Sugar().Errorf("failed load cert key from path %s %s", tlsCert, tlsKey)
		} else {
			tlsConfig = &tls.Config{
				Certificates: []tls.Certificate{cert},
			}
		}
	}

	var resolver *dns.Client
	var coreDnsAddr string
	if types.AgentConfig.AppDnsUpstream {
		resolver = &dns.Client{
			Net: "udp",
		}
		dnsServiceIPs, err := k8sObjManager.GetK8sObjManager().ListServicesDnsIP(context.Background())
		if err != nil {
			logger.Sugar().Errorf("failed get kube dns service,err: %v", err)
		}

		// If the ip address list is empty, the corresponding service does not exist or the selected label is incorrect
		if len(dnsServiceIPs) == 0 {
			logger.Sugar().Errorf("failed get kube dns service: %v", "the corresponding service does not exist or the selected label is incorrect")
			coreDnsAddr = ""
		} else {
			logger.Sugar().Infof("kube dns service %s ", dnsServiceIPs)
			coreDnsAddr = fmt.Sprintf("%s:53", dnsServiceIPs[0])
		}

	}
	handler := dns.HandlerFunc(func(w dns.ResponseWriter, r *dns.Msg) {
		qname := r.Question[0].Name
		e2e := strings.HasPrefix(qname, "netdns-e2e")
		m := new(dns.Msg)
		if e2e {
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
			task := strings.Split(qname, ".")[0]
			agentHttpServer.RequestCounts.AddOneCount(task)
		} else if types.AgentConfig.AppDnsUpstream && coreDnsAddr != "" {
			m, _, err = resolver.Exchange(r, coreDnsAddr)
			if err != nil {
				fmt.Println("Error forwarding DNS query:", err)
				return
			}
		}
		_ = w.WriteMsg(m)

	})

	if tlsConfig != nil {
		tlsServer := &dns.Server{
			Addr:      fmt.Sprintf(":%d", types.AgentConfig.AppDnsTcpTlsPort),
			Net:       "tcp-tls",
			TLSConfig: tlsConfig,
			Handler:   handler,
		}
		go func() {
			logger.Sugar().Infof("dns tcp-tls server listien %s", tlsServer.Addr)
			if err := tlsServer.ListenAndServe(); err != nil {
				logger.Sugar().Errorf("dns tcp-tls server , err: %v ", err)
			}
		}()
	}

	udpServer := &dns.Server{
		Addr:    fmt.Sprintf(":%d", types.AgentConfig.AppDnsUdpPort),
		Net:     "udp",
		Handler: handler,
	}

	go func() {
		logger.Sugar().Infof("dns udp server listien %s", udpServer.Addr)
		if err := udpServer.ListenAndServe(); err != nil {
			logger.Sugar().Errorf("dns udp server , err: %v ", err)
		}
	}()

	tcpServer := &dns.Server{
		Addr:    fmt.Sprintf(":%d", types.AgentConfig.AppDnsTcpPort),
		Net:     "tcp",
		Handler: handler,
	}
	go func() {
		logger.Sugar().Infof("dns tcp server listien %s", tcpServer.Addr)
		if err := tcpServer.ListenAndServe(); err != nil {
			logger.Sugar().Errorf("dns tcp server , err: %v ", err)
		}
	}()

}
