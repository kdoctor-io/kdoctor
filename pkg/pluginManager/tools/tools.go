// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0
//
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// function ValidataAppHttpHealthyHost base on https://github.com/golang/go/blob/master/src/net/ip.go  ParseIP function
//
// Changes:
// - add domain check

package tools

import (
	"fmt"
	crd "github.com/kdoctor-io/kdoctor/pkg/k8s/apis/kdoctor.io/v1beta1"
	"github.com/robfig/cron"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"net"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

var DomainRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]+(\.[a-zA-Z0-9_-]+)*\.[a-zA-Z]{2,}$`)

func ValidataCrdSchedule(plan *crd.SchedulePlan) error {

	if plan == nil {
		return fmt.Errorf("Schedule is empty ")
	}

	args := strings.Split(*plan.Schedule, " ")

	if len(args) == 2 {
		startAfterMinute, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("The format of the schedule is incorrect, it should be number ")
		}
		intervalMinute, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("The format of the schedule is incorrect, it should be number ")
		}
		if startAfterMinute < 0 {
			return fmt.Errorf("Schedule.StartAfterMinute %v must not be smaller than 0 ", startAfterMinute)
		}

		if intervalMinute < 1 {
			return fmt.Errorf("Schedule.IntervalMinute %v must not be smaller than 1 ", intervalMinute)
		}

		if int(plan.RoundTimeoutMinute) > intervalMinute {
			return fmt.Errorf("Schedule.RoundTimeoutMinute %v must not be bigger than Schedule.IntervalMinute %v ", plan.RoundTimeoutMinute, intervalMinute)
		}

	} else if len(args) == 5 {
		_, err := cron.ParseStandard(*plan.Schedule)
		if err != nil {
			return fmt.Errorf("Crontab configuration error,err: %v ", err)
		}

	} else {
		return fmt.Errorf("The format of the schedule is incorrect, it should be two or five ")
	}

	if plan.RoundTimeoutMinute < 1 {
		return fmt.Errorf("Schedule.RoundTimeoutMinute %v must not be smaller than 1 ", plan.RoundTimeoutMinute)
	}

	return nil
}

// ValidataAppHttpHealthyHost check host protocol,ipv4 and ipv6 addr
func ValidataAppHttpHealthyHost(r *crd.AppHttpHealthy) error {
	var ip string
	u, err := url.Parse(r.Spec.Target.Host)
	if err != nil {
		s := fmt.Sprintf("HttpAppHealthy %v, The IP address of the host is incorrect,err: %v ", r.Name, err)
		return apierrors.NewBadRequest(s)
	}
	ip = u.Hostname()
	for i := 0; i < len(ip); i++ {
		// ipv4 or domain
		if ip[i] == '.' {
			ipAddr := ip
			// if host contains port remove port
			if strings.Contains(ip, ":") {
				ipAddr = ip[:strings.Index(ip, ":")]
			}
			if DomainRegex.MatchString(ipAddr) {
				break
			}
			if net.ParseIP(ipAddr).To4() == nil {
				s := fmt.Sprintf("HttpAppHealthy %v, The IP address of the host is incorrect", r.Name)
				return apierrors.NewBadRequest(s)
			}
			break

			// ipv6
		} else if ip[i] == ':' {
			if r.Spec.Target.Host[strings.Index(r.Spec.Target.Host, ip)-1] != '[' {
				s := fmt.Sprintf("HttpAppHealthy %v, using bad/illegal format or missing URL,example: http://[ipv6]:port", r.Name)
				return apierrors.NewBadRequest(s)
			}
			if r.Spec.Target.Host[strings.Index(r.Spec.Target.Host, ip)+len(ip)] != ']' {
				s := fmt.Sprintf("HttpAppHealthy %v, using bad/illegal format or missing URL,example: http://[ipv6]:port", r.Name)
				return apierrors.NewBadRequest(s)
			}
			if net.ParseIP(ip).To16() == nil {
				s := fmt.Sprintf("HttpAppHealthy %v, The IP address of the host is incorrect", r.Name)
				return apierrors.NewBadRequest(s)
			}
			break
		}
	}

	return nil
}

func GetDefaultSchedule() (plan *crd.SchedulePlan) {
	s := "0 60"
	return &crd.SchedulePlan{
		RoundTimeoutMinute: 60,
		Schedule:           &s,
		RoundNumber:        1,
	}
}

func GetDefaultNetSuccessCondition() (plan *crd.NetSuccessCondition) {
	n := float64(1)
	return &crd.NetSuccessCondition{
		SuccessRate: &n,
	}
}
