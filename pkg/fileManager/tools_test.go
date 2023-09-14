// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0
package fileManager_test

import (
	"github.com/kdoctor-io/kdoctor/pkg/fileManager"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("test fileManager tool", Label("fileManager"), func() {

	It("test basic", func() {
		filePath := "/tmp/_loggertest/a.txt"

		fileManager.DefaultFileWriter(100, 0, 0)
		wr := fileManager.NewFileWriter(filePath)
		GinkgoWriter.Printf("succeed to new write for %v", filePath)
		defer wr.Close()

		data := []byte("test data\n dsf\n")
		_, t := wr.Write(data)
		Expect(t).NotTo(HaveOccurred())

	})

})
