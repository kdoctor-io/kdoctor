// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/api/errors"

	k8sObjManager "github.com/kdoctor-io/kdoctor/pkg/k8ObjManager"
	"github.com/kdoctor-io/kdoctor/pkg/types"
	"github.com/kdoctor-io/kdoctor/pkg/utils"
)

var (
	TlsCertPath = "/tmp/cert.crt"
	TlsKeyPath  = "/tmp/key.crt"
	CaCertPath  = "/tmp/ca.crt"
)

func GenServerCert(logger *zap.Logger) {
	checkServiceReady(logger)

	// get svc domain and ip
	alternateIP := []net.IP{}
	alternateDNS := []string{}
	servicePortName := "http"

	if types.AgentConfig.Configmap.EnableIPv4 {
		serviceIPv4, err := k8sObjManager.GetK8sObjManager().GetServiceAccessUrl(context.Background(), types.AgentConfig.ServiceV4Name, types.AgentConfig.PodNamespace, servicePortName)
		if err != nil {
			logger.Sugar().Fatalf("failed to get kdoctor ipv4 service %s/%s, reason=%v ", types.AgentConfig.PodNamespace, types.AgentConfig.ServiceV4Name, err)
		}
		logger.Sugar().Debugf("get ipv4 serviceAccessurl %v", serviceIPv4)
		// ipv4 ip
		for _, v := range serviceIPv4.ClusterIPUrl {
			host, _, err := net.SplitHostPort(v)
			if err != nil {
				logger.Sugar().Errorf("ip addr %s split host port err,reason: %v ", v, err)
				continue
			}
			alternateIP = append(alternateIP, net.ParseIP(host))
		}

		for _, v := range serviceIPv4.LoadBalancerUrl {
			host, _, err := net.SplitHostPort(v)
			if err != nil {
				logger.Sugar().Errorf("ip addr %s split host port err,reason: %v ", v, err)
				continue
			}
			alternateIP = append(alternateIP, net.ParseIP(host))
		}

		// ipv4 dns
		alternateDNS = append(alternateDNS, types.AgentConfig.ServiceV4Name)
		domain := fmt.Sprintf("%s.%s", types.AgentConfig.ServiceV4Name, types.AgentConfig.PodNamespace)
		alternateDNS = append(alternateDNS, domain)
		alternateDNS = append(alternateDNS, fmt.Sprintf("%s.svc", domain))
		alternateDNS = append(alternateDNS, fmt.Sprintf("%s.svc.cluster.local", domain))
	}

	if types.AgentConfig.Configmap.EnableIPv6 {
		serviceIPv6, err := k8sObjManager.GetK8sObjManager().GetServiceAccessUrl(context.Background(), types.AgentConfig.ServiceV6Name, types.AgentConfig.PodNamespace, servicePortName)
		if err != nil {
			logger.Sugar().Fatalf("failed to get kdoctor ipv6 service %s/%s, reason=%v ", types.AgentConfig.PodNamespace, types.AgentConfig.ServiceV6Name, err)
		}
		// ipv6 ip
		logger.Sugar().Debugf("get ipv6 serviceAccessurl %v", serviceIPv6)
		for _, v := range serviceIPv6.ClusterIPUrl {
			p := strings.LastIndex(v, ":")
			host := net.ParseIP(v[:p])
			if len(host) == 0 {
				logger.Sugar().Errorf("parse ip %s err", v[:p])
				continue
			}
			alternateIP = append(alternateIP, host)
		}
		for _, v := range serviceIPv6.LoadBalancerUrl {
			p := strings.LastIndex(v, ":")
			host := net.ParseIP(v[:p])
			if len(host) == 0 {
				logger.Sugar().Errorf("parse ip %s err", v[:p])
				continue
			}
			alternateIP = append(alternateIP, host)
		}

		// ipv6 dns
		alternateDNS = append(alternateDNS, types.AgentConfig.ServiceV6Name)
		domain := fmt.Sprintf("%s.%s", types.AgentConfig.ServiceV6Name, types.AgentConfig.PodNamespace)
		alternateDNS = append(alternateDNS, domain)
		alternateDNS = append(alternateDNS, fmt.Sprintf("%s.svc", domain))
		alternateDNS = append(alternateDNS, fmt.Sprintf("%s.svc.cluster.local", domain))
	}

	alternateDNS = append(alternateDNS, types.AgentConfig.PodName)
	logger.Sugar().Debugf("alternate ip for tls cert:  %v", alternateIP)
	logger.Sugar().Debugf("alternate dns for tls cert:  %v", alternateDNS)
	// generate self-signed certificates
	if e := utils.NewServerCertKeyForLocalNode(alternateDNS, alternateIP, types.AgentConfig.TlsCaCertPath, types.AgentConfig.TlsCaKeyPath, TlsCertPath, TlsKeyPath, CaCertPath); e != nil {
		logger.Sugar().Fatalf("failed to generate certiface, error=%v", e)
	}
}

func checkServiceReady(logger *zap.Logger) {
	ctx := context.TODO()

	if types.AgentConfig.Configmap.EnableIPv4 {
		for {
			_, err := k8sObjManager.GetK8sObjManager().GetService(ctx, types.AgentConfig.ServiceV4Name, types.AgentConfig.PodNamespace)
			if nil != err {
				if errors.IsNotFound(err) {
					logger.Sugar().Errorf("agent runtime IPv4 service %s/%s not exists, wait for controller to create it", types.AgentConfig.PodNamespace, types.AgentConfig.ServiceV4Name)
					time.Sleep(time.Second)
					continue
				}
				logger.Sugar().Errorf("failed to get agent runtime IPv4 service %s/%s, error: %v", types.AgentConfig.PodNamespace, types.AgentConfig.ServiceV4Name, err)
			} else {
				break
			}
		}
	}

	if types.AgentConfig.Configmap.EnableIPv6 {
		for {
			_, err := k8sObjManager.GetK8sObjManager().GetService(ctx, types.AgentConfig.ServiceV6Name, types.AgentConfig.PodNamespace)
			if nil != err {
				if errors.IsNotFound(err) {
					logger.Sugar().Errorf("agent runtime IPv6 service %s/%s not exists, wait for controller to create it", types.AgentConfig.PodNamespace, types.AgentConfig.ServiceV6Name)
					time.Sleep(time.Second)
					continue
				}
				logger.Sugar().Errorf("failed to get runtime agent IPv6 service %s/%s, error: %v", types.AgentConfig.PodNamespace, types.AgentConfig.ServiceV6Name, err)
			} else {
				break
			}
		}
	}
}
