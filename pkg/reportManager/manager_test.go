// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package reportManager

import (
	"context"
	"fmt"
	"github.com/agiledragon/gomonkey/v2"
	"github.com/golang/mock/gomock"
	k8sObjManager "github.com/kdoctor-io/kdoctor/pkg/k8ObjManager"
	"github.com/kdoctor-io/kdoctor/pkg/k8s/apis/kdoctor.io/v1beta1"
	"github.com/kdoctor-io/kdoctor/pkg/logger"
	"github.com/kdoctor-io/kdoctor/pkg/scheduler"
	"github.com/kdoctor-io/kdoctor/pkg/types"
	"github.com/kdoctor-io/kdoctor/pkg/utils"
	. "github.com/onsi/ginkgo/v2"
	"time"
)

var _ = Describe("unit test", Label("unit test"), func() {
	var reportDir string
	BeforeEach(func() {
		reportDir = fmt.Sprintf("/tmp/_FM_%d", time.Now().Nanosecond())
	})

	It("init InitReportManager ", Label("reportManager"), func() {

		ctx := context.Background()
		// mock
		patch := gomonkey.ApplyFuncReturn(utils.GetFileList, []string{"test-report-local-file"}, nil)
		defer patch.Reset()

		mockPodIP := k8sObjManager.PodIps{"testPod": []k8sObjManager.IPs{{
			InterfaceName: "eth0",
			IPv4:          "127.0.0.1",
		}}}

		k8sManager.EXPECT().ListDeploymentPodIPs(gomock.Eq(ctx), gomock.Eq("deployment"), gomock.Eq("default")).Return(mockPodIP, nil).AnyTimes()
		k8sManager.EXPECT().ListDaemonsetPodIPs(gomock.Eq(ctx), gomock.Eq("daemonset"), gomock.Eq("default")).Return(mockPodIP, nil).AnyTimes()

		grpcClient.EXPECT().GetFileList(gomock.Eq(ctx), gomock.Eq("127.0.0.1"), gomock.Eq("/report")).Return([]string{"test-report-local-file"}, nil).AnyTimes()

		log := logger.NewStdoutLogger("debug", "reportManager Test")
		dbMap := make(map[string]scheduler.DB, 3)

		appDB := scheduler.NewDB(5, log)
		reachDB := scheduler.NewDB(5, log)
		dnsDB := scheduler.NewDB(5, log)

		dbMap[types.KindNameAppHttpHealthy] = appDB
		dbMap[types.KindNameNetReach] = reachDB
		dbMap[types.KindNameNetdns] = dnsDB

		InitReportManager(log, reportDir, time.Second*5, dbMap)

		_ = appDB.Apply(scheduler.BuildItem(v1beta1.TaskResource{
			RuntimeName:   "test-app",
			RuntimeType:   types.KindDeployment,
			RuntimeStatus: v1beta1.RuntimeCreated,
		}, types.KindNameAppHttpHealthy, "test-app", nil))

		_ = reachDB.Apply(scheduler.BuildItem(v1beta1.TaskResource{
			RuntimeName:   "test-reach",
			RuntimeType:   types.KindDeployment,
			RuntimeStatus: v1beta1.RuntimeCreated,
		}, types.KindNameNetReach, "test-reach", nil))

		_ = dnsDB.Apply(scheduler.BuildItem(v1beta1.TaskResource{
			RuntimeName:   "test-dns",
			RuntimeType:   types.KindDeployment,
			RuntimeStatus: v1beta1.RuntimeCreated,
		}, types.KindNameNetdns, "test-dns", nil))

		TriggerSyncReport("test-app.1")
		TriggerSyncReport("test-reach.1")
		TriggerSyncReport("test-dns.1")
		time.Sleep(time.Second * 20)
	})

})
