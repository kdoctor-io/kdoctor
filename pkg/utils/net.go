// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"fmt"
	"github.com/pkg/errors"
	"net"
	"strings"
)

func ListHostAllInterfaces() ([]string, error) {
	list, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	r := []string{}
	for _, v := range list {
		r = append(r, v.Name)
	}
	if len(r) == 0 {
		return nil, nil
	}
	return r, nil
}

// 1.1.1.1
func CheckIPv4Format(ip string) bool {
	result := net.ParseIP(ip)
	if result == nil {
		return false
	}
	if result.To4() == nil {
		return false
	}
	return true
}
func CheckIPv6Format(ip string) bool {
	result := net.ParseIP(ip)
	if result == nil {
		return false
	}
	if result.To16() == nil {
		return false
	}
	return true
}

func GetInterfaceUnicastAddrByName(name string) (ipv4MaskList, ipv6MaskList []string, err error) {
	ipv4MaskList = []string{}
	ipv6MaskList = []string{}
	result, err := net.InterfaceByName(name)
	if err != nil {
		return nil, nil, fmt.Errorf("no host interface with name=%v ", name)
	}
	// returns a list of unicast interface addresses
	list, e := result.Addrs()
	if e != nil {
		return nil, nil, e
	}
	for _, v := range list {
		m := v.String()
		if CheckIPv4Format(strings.Split(m, "/")[0]) {
			ipv4MaskList = append(ipv4MaskList, m)
		} else {
			ipv6MaskList = append(ipv6MaskList, m)
		}
	}
	return
}

func GetAllInterfaceUnicastAddr() (ipv4MaskList, ipv6MaskList []string, err error) {
	ipv4MaskList = []string{}
	ipv6MaskList = []string{}

	re, err := ListHostAllInterfaces()
	if err != nil {
		return nil, nil, err
	}

	for _, name := range re {
		m, n, err := GetInterfaceUnicastAddrByName(name)
		if err != nil {
			return nil, nil, err
		}
		if len(m) > 0 {
			ipv4MaskList = append(ipv4MaskList, m...)
		}
		if len(n) > 0 {
			ipv6MaskList = append(ipv6MaskList, n...)
		}
	}
	return
}

func GetAllInterfaceUnicastAddrWithoutMask() (ipv4List, ipv6List []net.IP, err error) {
	ipv4List = []net.IP{}
	m, n, e := GetAllInterfaceUnicastAddr()
	if e != nil {
		return nil, nil, e
	}
	for _, t := range m {
		s := strings.Split(t, "/")[0]
		q := net.ParseIP(s)
		if len(q) == 0 {
			return nil, nil, errors.Errorf("fail to parse %v", s)
		}
		ipv4List = append(ipv4List, q)
	}
	for _, t := range n {
		s := strings.Split(t, "/")[0]
		q := net.ParseIP(s)
		if len(q) == 0 {
			return nil, nil, errors.Errorf("fail to parse %v", s)
		}
		ipv6List = append(ipv6List, q)
	}
	return ipv4List, ipv6List, nil
}
