// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package reportManager

import (
	"context"
	"fmt"
	"github.com/agiledragon/gomonkey/v2"
	"github.com/golang/mock/gomock"
	"github.com/kdoctor-io/kdoctor/pkg/logger"
	"github.com/kdoctor-io/kdoctor/pkg/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/client-go/util/workqueue"
	"time"
)

var _ = Describe("unit test", Label("unit test "), func() {
	var reportDir string
	ctx := context.Background()
	BeforeEach(func() {
		reportDir = fmt.Sprintf("/tmp/_FM_%d", time.Now().Nanosecond())
	})

	It("test getMissRemoteReport", Label("reportManager worker"), func() {

		name1 := "Nethttp_test-agent_round1_kdoctor-worker_2022-12-21T11:47:08Z"
		name2 := "Nethttp_test-agent_round1_kdoctor-control-plane_2022-12-21T11:45:08Z"
		remoteFileList := []string{name1, name2}
		localFileList := []string{name1}

		missFile := GetMissRemoteReport(remoteFileList, localFileList)
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

		missFile := GetMissRemoteReport(remoteFileList, localFileList)
		fmt.Printf("%v", missFile)

	})

	It("syncReportFromOneAgent ", Label("reportManager worker"), func() {
		types.ControllerConfig.DirPathAgentReport = "/report"
		// mock
		patch := gomonkey.ApplyFuncReturn(GetMissRemoteReport, []string{"Nethttp_test-agent_round1_kdoctor-worker_2022-12-21T12:18:20Z"})
		defer patch.Reset()

		grpcClient.EXPECT().GetFileList(gomock.Eq(ctx), gomock.Any(), gomock.Any()).Return([]string{"Nethttp_test-agent_round1_kdoctor-control-plane_2022-12-21T12:18:20Z"}, nil).AnyTimes()
		grpcClient.EXPECT().SaveRemoteFileToLocal(gomock.Eq(ctx), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

		log := logger.NewStdoutLogger("debug", "reportManager Test")

		rm := &reportManager{
			logger:          log,
			reportDir:       reportDir,
			collectInterval: time.Second * 5,
			queue:           workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "reportManager"),
		}

		rm.syncReportFromOneAgent(ctx, log, grpcClient, []string{"Nethttp_test-agent_round1_kdoctor-worker_2022-12-21T12:18:20Z"}, "test", "127.0.0.1")

	})

})
