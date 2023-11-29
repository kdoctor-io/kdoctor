// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package reportManager

import (
	"context"
	"github.com/kdoctor-io/kdoctor/pkg/scheduler"
	"github.com/kdoctor-io/kdoctor/pkg/types"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/util/workqueue"
	"time"
)

const (
	queueMaxRetries = 100
)

type reportManager struct {
	logger                  *zap.Logger
	reportDir               string
	collectInterval         time.Duration
	queue                   workqueue.RateLimitingInterface
	appHttpHealthyRuntimeDB scheduler.DB
	netReachRuntimeDB       scheduler.DB
	netDNSRuntimeDB         scheduler.DB
}

var globalReportManager *reportManager

func InitReportManager(logger *zap.Logger, reportDir string, collectInterval time.Duration, db map[string]scheduler.DB) {
	if globalReportManager != nil {
		return
	}

	globalReportManager = &reportManager{
		logger:          logger,
		reportDir:       reportDir,
		collectInterval: collectInterval,
		queue:           workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "reportManager"),
	}

	for k, v := range db {
		switch k {
		case types.KindNameAppHttpHealthy:
			globalReportManager.appHttpHealthyRuntimeDB = v
		case types.KindNameNetReach:
			globalReportManager.netReachRuntimeDB = v
		case types.KindNameNetdns:
			globalReportManager.netDNSRuntimeDB = v
		}
	}

	go globalReportManager.runWorker()
}

func (s *reportManager) runWorker() {

	// TODO: wait for all agent grpc is ready
	s.logger.Info("waiting for all agent grpc server ready")

	s.logger.Info("all agent grpc server ready, start worker")

	//
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	defer s.queue.ShutDown()

	// please do not run more than one worker, or else it races to write reports
	go wait.UntilWithContext(ctx, s.worker, time.Second)

	select {}
}

func TriggerSyncReport(tgt string) {
	if globalReportManager != nil {
		globalReportManager.logger.Sugar().Debugf("trigger to sync agent report from source %v", tgt)
		// s.queue.AddRateLimited(triggerName)
		globalReportManager.queue.AddAfter(tgt, 10*time.Second)
	}
}

// --------------

func (s *reportManager) worker(ctx context.Context) {
	for s.processNextWorkItem(ctx) {
	}
}

func (s *reportManager) processNextWorkItem(ctx context.Context) bool {
	key, quit := s.queue.Get()
	if quit {
		return false
	}
	defer s.queue.Done(key)

	err := s.syncHandler(ctx, key.(string))
	if err == nil {
		s.queue.Forget(key)
	} else {
		s.logger.Sugar().Warnf("worker failed , error=%v", err)
		if s.queue.NumRequeues(key) < queueMaxRetries {
			s.queue.AddRateLimited(key)
		}
	}
	// handle nex item
	return true
}
