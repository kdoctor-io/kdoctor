// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0
package loadDns_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestLoadDns(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "load request Suite")
}

var _ = BeforeSuite(func() {
	// nothing to do
})
