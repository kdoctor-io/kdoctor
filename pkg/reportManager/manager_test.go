// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package reportManager_test

import (
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/kdoctor-io/kdoctor/pkg/reportManager"
)

var _ = Describe("unit test", Label("unit test"), func() {

	It("test getMissRemoteReport", Pending, func() {

		name1 := "Nethttp_test-agent_round1_kdoctor-worker_2022-12-21T11:47:08Z"
		name2 := "Nethttp_test-agent_round1_kdoctor-control-plane_2022-12-21T11:45:08Z"
		remoteFileList := []string{name1, name2}
		localFileList := []string{name1}

		missFile := reportManager.GetMissRemoteReport(remoteFileList, localFileList)
		Expect(len(missFile)).To(Equal(1))
		Expect(missFile).To(ConsistOf([]string{name2}))

	})

	It("test2 getMissRemoteReport", func() {

		remoteFileList := []string{
			"Nethttp_test-agent_round1_kdoctor-worker_2022-11-21T12:24:08Z",
			"Nethttp_test-agent_round2_kdoctor-worker_2022-11-21T12:26:07Z",
		}

		localFileList := []string{
			"Nethttp_test-agent_round1_kdoctor-control-plane_2022-12-21T12:18:20Z",
			"Nethttp_test-agent_round1_kdoctor-worker_2022-12-21T12:18:20Z",
			"Nethttp_test-agent_round1_summary_2022-12-21T12:18:05Z",
			"Nethttp_test-agent_round2_summary_2022-12-21T12:20:05Z",
		}

		missFile := reportManager.GetMissRemoteReport(remoteFileList, localFileList)
		fmt.Printf("%v", missFile)

	})

})
