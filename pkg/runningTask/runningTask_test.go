// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package runningTask

import (
	"github.com/kdoctor-io/kdoctor/pkg/k8s/apis/system/v1beta1"
	"github.com/kdoctor-io/kdoctor/pkg/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("runningTask unit test", Label("runningTask"), func() {

	It("running task", func() {
		var load v1beta1.TotalRunningLoad
		rt := InitRunningTask()
		rt.SetTask(Task{Kind: types.KindNameAppHttpHealthy, Qps: 10, Name: "test1"})
		rt.SetTask(Task{Kind: types.KindNameNetdns, Qps: 20, Name: "test2"})
		rt.SetTask(Task{Kind: types.KindNameNetReach, Qps: 30, Name: "test3"})
		load = rt.QpsStats()

		Expect(load.AppHttpHealthyQPS).To(Equal(int64(10)), "appHttpHealthy qps not equal")
		Expect(load.NetReachQPS).To(Equal(int64(30)), "netReach qps not equal")
		Expect(load.NetDnsQPS).To(Equal(int64(20)), "netDns qps not equal")

		rt.DeleteTask("test1")
		load = rt.QpsStats()
		Expect(load.AppHttpHealthyQPS).To(Equal(int64(0)), "appHttpHealthy qps not equal")
		Expect(load.NetReachQPS).To(Equal(int64(30)), "netReach qps not equal")
		Expect(load.NetDnsQPS).To(Equal(int64(20)), "netDns qps not equal")

	})

})
