// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package runtime

import (
	"context"
	"github.com/kdoctor-io/kdoctor/pkg/logger"
	"github.com/kdoctor-io/kdoctor/pkg/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("runtime unit test", Label("runtime"), func() {

	It("runtime task daemonSet", func() {
		name := "testDaemonSet"
		namespace := "testDaemonSet"
		ctx := context.Background()
		runtimeDaemonSet := NewDaemonSetRuntime(c, c, namespace, name, logger.NewStdoutLogger("debug", "runtime daemonSet"))

		ds := types.DaemonsetTempl.DeepCopy()
		ds.SetName(name)
		ds.SetNamespace(namespace)

		err := c.Create(ctx, ds)
		Expect(err).To(BeNil(), "create daemonSet failed")

		ready := runtimeDaemonSet.IsReady(ctx)
		Expect(ready).To(BeTrue(), "not ready daemonset ")
		err = runtimeDaemonSet.Delete(ctx)
		Expect(err).To(BeNil(), "delete daemonSet failed")

		runtimeDaemonSetNotFound := NewDaemonSetRuntime(c, c, "not-found", "not-found", logger.NewStdoutLogger("debug", "runtime daemonSet"))
		readyNotFound := runtimeDaemonSetNotFound.IsReady(ctx)
		Expect(readyNotFound).To(BeFalse(), "not ready daemonset ")
		err = runtimeDaemonSetNotFound.Delete(ctx)
		Expect(err).To(BeNil(), "not found daemonSet")
		_, err = FindRuntime(c, c, types.KindDaemonSet, namespace, logger.NewStdoutLogger("debug", "runtime daemonSet"))
		Expect(err).To(BeNil(), "get runtime")

	})

	It("runtime task deployment", func() {
		name := "testDeployment"
		namespace := "testDeployment"
		ctx := context.Background()
		runtimeDeploy := NewDeploymentRuntime(c, c, namespace, name, logger.NewStdoutLogger("debug", "runtime deploy"))

		deploy := types.DeploymentTempl.DeepCopy()
		deploy.SetName(name)
		deploy.SetNamespace(namespace)

		err := c.Create(ctx, deploy)
		Expect(err).To(BeNil(), "create deployment failed")

		runtimeDeploy.IsReady(ctx)
		err = runtimeDeploy.Delete(ctx)
		Expect(err).To(BeNil(), "delete deployment failed")

		runtimeDeploymentNotFound := NewDeploymentRuntime(c, c, "not-found", "not-found", logger.NewStdoutLogger("debug", "runtime daemonSet"))
		readyNotFound := runtimeDeploymentNotFound.IsReady(ctx)
		Expect(readyNotFound).To(BeFalse(), "not ready deployment ")
		err = runtimeDeploymentNotFound.Delete(ctx)
		Expect(err).To(BeNil(), "not found deployment")
		_, err = FindRuntime(c, c, types.KindDeployment, namespace, logger.NewStdoutLogger("debug", "runtime daemonSet"))
		Expect(err).To(BeNil(), "get runtime")
	})

	It("FindRuntime test ", Label("runtime"), func() {
		_, err := FindRuntime(c, c, "test", "test", logger.NewStdoutLogger("debug", "runtime"))
		Expect(err).NotTo(BeNil(), "unrecognized runtime type")
	})
})
