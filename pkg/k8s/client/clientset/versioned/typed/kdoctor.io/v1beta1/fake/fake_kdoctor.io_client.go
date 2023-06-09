// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	v1beta1 "github.com/kdoctor-io/kdoctor/pkg/k8s/client/clientset/versioned/typed/kdoctor.io/v1beta1"
	rest "k8s.io/client-go/rest"
	testing "k8s.io/client-go/testing"
)

type FakeKdoctorV1beta1 struct {
	*testing.Fake
}

func (c *FakeKdoctorV1beta1) AppHttpHealthies() v1beta1.AppHttpHealthyInterface {
	return &FakeAppHttpHealthies{c}
}

func (c *FakeKdoctorV1beta1) NetReaches() v1beta1.NetReachInterface {
	return &FakeNetReaches{c}
}

func (c *FakeKdoctorV1beta1) Netdnses() v1beta1.NetdnsInterface {
	return &FakeNetdnses{c}
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *FakeKdoctorV1beta1) RESTClient() rest.Interface {
	var ret *rest.RESTClient
	return ret
}
