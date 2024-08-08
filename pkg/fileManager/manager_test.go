// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package fileManager_test

import (
	"fmt"
	"os"
	"path"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"time"

	"github.com/kdoctor-io/kdoctor/pkg/fileManager"
	"github.com/kdoctor-io/kdoctor/pkg/logger"
)

var _ = Describe("test fileManager", Label("fileManager"), func() {

	var reportDir string

	BeforeEach(func() {
		reportDir = fmt.Sprintf("/tmp/_FM_%d", time.Now().Nanosecond())
	})

	It("test basic", func() {
		log := logger.NewStdoutLogger("debug", "test")
		cleanInterval := 2 * time.Second
		f, e := fileManager.NewManager(log, reportDir, cleanInterval)
		Expect(e).NotTo(HaveOccurred(), "failed to NewManager, error=%v", e)

		// --write
		kindName := "kindTom"
		taskName := "taskFire"
		nodeName := "worker1"
		roundNumber := 10
		endTime := time.Now().Add(10 * time.Second)
		data := []byte("line1 \n line2\n")
		e = f.WriteTaskFile(kindName, taskName, roundNumber, nodeName, endTime, data)
		Expect(e).NotTo(HaveOccurred(), "failed to write task file %v", e)

		// ---- check existence
		time.Sleep(5 * time.Second)
		expectedFileName := fileManager.GenerateTaskFileName(kindName, taskName, roundNumber, nodeName, endTime)
		expectedFilePath := path.Join(reportDir, expectedFileName)
		GinkgoWriter.Printf("expect file %v \n", expectedFilePath)

		filelist, e := f.GetTaskAllFile(kindName, taskName)
		Expect(e).NotTo(HaveOccurred(), "failed to read directory %s, error=%v", reportDir, e)
		Expect(len(filelist)).To(Equal(1))
		Expect(filelist).To(ConsistOf([]string{expectedFilePath}))
		Expect(f.CheckTaskFileExisted(kindName, taskName, roundNumber)).To(BeTrue())
		// read data
		readdata, er := os.ReadFile(expectedFilePath)
		Expect(er).NotTo(HaveOccurred(), "failed to read file %s, error=%v", expectedFilePath, er)
		GinkgoWriter.Printf("read data: %v \n", string(readdata))

		// ---- check deletion
		time.Sleep(10 * time.Second)
		filelist, e = f.GetTaskAllFile(kindName, taskName)
		Expect(e).NotTo(HaveOccurred(), "failed to read directory , error=%v", e)
		Expect(len(filelist)).To(Equal(0))

	})
})
