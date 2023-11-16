// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0
package reportManager

import (
	"github.com/golang/mock/gomock"
	grpcManager_mock "github.com/kdoctor-io/kdoctor/pkg/grpcManager/mock"
	k8sObjManager_mock "github.com/kdoctor-io/kdoctor/pkg/k8ObjManager/mock"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var mockCtrl *gomock.Controller
var k8sManager *k8sObjManager_mock.MockK8sObjManager
var grpcClient *grpcManager_mock.MockGrpcClientManager

func TestReportManager(t *testing.T) {
	mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()
	RegisterFailHandler(Fail)
	RunSpecs(t, "reportManager Suite")
}

var _ = BeforeSuite(func() {
	// nothing to do
	k8sManager = k8sObjManager_mock.NewMockK8sObjManager(mockCtrl)
	grpcClient = grpcManager_mock.NewMockGrpcClientManager(mockCtrl)

})
