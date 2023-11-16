// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package scheduler

import (
	"github.com/kdoctor-io/kdoctor/pkg/k8s/apis/kdoctor.io/v1beta1"
	"github.com/kdoctor-io/kdoctor/pkg/logger"
	"github.com/kdoctor-io/kdoctor/pkg/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
)

var _ = Describe("cacheDB unit test", Label("cacheDB"), func() {

	It("cacheDB", Label("cacheDB"), func() {
		db := NewDB(30, logger.NewStdoutLogger("debug", "cacheDB"))
		taskName := "test"
		runtimeName := TaskRuntimeName(types.KindNameAppHttpHealthy, taskName)
		serviceV4Name := TaskRuntimeServiceName(runtimeName, corev1.IPv4Protocol)
		serviceV6Name := TaskRuntimeServiceName(runtimeName, corev1.IPv6Protocol)
		item := BuildItem(v1beta1.TaskResource{
			RuntimeName:   runtimeName,
			RuntimeType:   types.KindDeployment,
			ServiceNameV4: &serviceV4Name,
			ServiceNameV6: &serviceV6Name,
			RuntimeStatus: v1beta1.RuntimeCreated,
		}, types.KindNameAppHttpHealthy, taskName, nil)

		err := db.Apply(item)
		Expect(err).To(BeNil(), "apply db")

		runtimeList := db.List()
		Expect(len(runtimeList)).To(Equal(1), "runtime list")

		runtimeGet, err := db.Get(taskName)
		Expect(err).To(BeNil(), "get runtime")
		Expect(runtimeGet.RuntimeName).To(Equal(item.RuntimeName), "get runtime")
		Expect(runtimeGet.RuntimeKind).To(Equal(item.RuntimeKind), "get runtime")
		Expect(runtimeGet.RuntimeKey).To(Equal(item.RuntimeKey), "get runtime")
		Expect(runtimeGet.RuntimeStatus).To(Equal(item.RuntimeStatus), "get runtime")
		Expect(runtimeGet.ServiceNameV4).To(Equal(item.ServiceNameV4), "get runtime")
		Expect(runtimeGet.ServiceNameV6).To(Equal(item.ServiceNameV6), "get runtime")
		Expect(runtimeGet.RuntimeDeletionTime).To(Equal(item.RuntimeDeletionTime), "get runtime")

		Expect(len(db.List())).To(Equal(1), "delete item")

		taskName2 := "test2"
		runtimeName2 := TaskRuntimeName(types.KindNameAppHttpHealthy, taskName2)
		serviceV4Name2 := TaskRuntimeServiceName(runtimeName2, corev1.IPv4Protocol)
		serviceV6Name2 := TaskRuntimeServiceName(runtimeName2, corev1.IPv6Protocol)
		item2 := BuildItem(v1beta1.TaskResource{
			RuntimeName:   runtimeName2,
			RuntimeType:   types.KindDeployment,
			ServiceNameV4: &serviceV4Name2,
			ServiceNameV6: &serviceV6Name2,
			RuntimeStatus: v1beta1.RuntimeCreated,
		}, types.KindNameAppHttpHealthy, taskName2, nil)

		err = db.Apply(item2)
		Expect(err).To(BeNil(), "apply other item")
		db.Delete(item)
		db.Delete(item)
		Expect(len(db.List())).To(Equal(1), "delete item")

		item3 := BuildItem(v1beta1.TaskResource{
			RuntimeName:   runtimeName2,
			RuntimeType:   types.KindDeployment,
			ServiceNameV4: &serviceV4Name2,
			ServiceNameV6: &serviceV6Name2,
			RuntimeStatus: v1beta1.RuntimeCreated,
		}, types.KindNameNetReach, taskName2, nil)
		err = db.Apply(item3)
		Expect(err).To(BeNil(), "apply other item")
		Expect(len(db.List())).To(Equal(1), "len item")

		_ = NewDB(0, logger.NewStdoutLogger("debug", "cacheDB"))
		db3 := NewDB(1, logger.NewStdoutLogger("debug", "cacheDB"))
		err = db3.Apply(item2)
		Expect(err).To(BeNil(), "apply other item")
		err = db3.Apply(item)
		Expect(err).NotTo(BeNil(), "apply other item")

		_, err = db3.Get(taskName)
		Expect(err).NotTo(BeNil(), "get not exist task")
	})

})
