// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package k8sObjManager

import (
	"context"
	"fmt"
	"github.com/kdoctor-io/kdoctor/pkg/types"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"net"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (nm *k8sObjManager) GetService(ctx context.Context, name, namespace string) (*corev1.Service, error) {

	d := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}
	key := client.ObjectKeyFromObject(d)
	if e := nm.client.Get(ctx, key, d); e != nil {
		return nil, fmt.Errorf("failed to get service %v/%v, reason=%v", namespace, name, e)
	}
	return d, nil
}

type ServiceAccessUrl struct {
	NodePort        int
	ClusterIPUrl    []string
	LoadBalancerUrl []string
}

func (nm *k8sObjManager) GetServiceAccessUrl(ctx context.Context, name, namespace string, portName string) (*ServiceAccessUrl, error) {
	s, e := nm.GetService(ctx, name, namespace)
	if e != nil || s == nil {
		return nil, e
	}

	r := &ServiceAccessUrl{
		ClusterIPUrl:    []string{},
		LoadBalancerUrl: []string{},
	}

	// get port index
	var portIndex int
	if len(portName) == 0 {
		t := len(s.Spec.Ports)
		if t > 1 {
			return nil, fmt.Errorf("variable portName is empty, but service has multiple ports: %v", s.Spec.Ports)
		}
		if t == 0 {
			return nil, fmt.Errorf("service has empty ports ")
		}
		portIndex = 0
	} else {
		found := false
		for l, v := range s.Spec.Ports {
			if v.Name == portName {
				portIndex = l
				found = true
				break
			}
		}
		if !found {
			return nil, fmt.Errorf("service has no port named %s", portName)
		}
	}

	// get clusterip url
	if len(s.Spec.ClusterIPs) > 0 {
		for _, v := range s.Spec.ClusterIPs {
			var t string
			if net.ParseIP(v).To4() == nil {
				t = fmt.Sprintf("[%s]:%v", v, s.Spec.Ports[portIndex].Port)
			} else {
				t = fmt.Sprintf("%s:%v", v, s.Spec.Ports[portIndex].Port)
			}
			r.ClusterIPUrl = append(r.ClusterIPUrl, t)
		}
	} else if len(s.Spec.ClusterIP) > 0 {
		var t string
		if net.ParseIP(s.Spec.ClusterIP).To4() == nil {
			t = fmt.Sprintf("[%s]:%v", s.Spec.ClusterIP, s.Spec.Ports[portIndex].Port)
		} else {
			t = fmt.Sprintf("%s:%v", s.Spec.ClusterIP, s.Spec.Ports[portIndex].Port)
		}
		r.ClusterIPUrl = append(r.ClusterIPUrl, t)
	}

	// get nodePort
	r.NodePort = int(s.Spec.Ports[portIndex].NodePort)

	// get loadbalancer url
	if len(s.Status.LoadBalancer.Ingress) > 0 {
		for _, v := range s.Status.LoadBalancer.Ingress {
			var t string
			if net.ParseIP(s.Spec.ClusterIP).To4() == nil {
				t = fmt.Sprintf("[%s]:%v", v.IP, s.Spec.Ports[portIndex].Port)
			} else {
				t = fmt.Sprintf("%s:%v", v.IP, s.Spec.Ports[portIndex].Port)
			}
			r.LoadBalancerUrl = append(r.LoadBalancerUrl, t)
		}
	}

	return r, nil
}

func (nm *k8sObjManager) ListServicesDnsIP(ctx context.Context) ([]string, error) {
	serviceList := new(corev1.ServiceList)
	var err error
	var result []string

	label := map[string]string{
		types.AgentConfig.DnsServiceSelectLabelKey: types.AgentConfig.DnsServiceSelectLabelValue,
	}

	ListOption := &client.ListOptions{
		Namespace:     "kube-system",
		LabelSelector: labels.SelectorFromSet(label),
	}

	if err = nm.client.List(ctx, serviceList, ListOption); err != nil {
		return nil, err
	}

	for _, v := range serviceList.Items {
		result = append(result, v.Spec.ClusterIPs...)
	}

	return result, nil
}
