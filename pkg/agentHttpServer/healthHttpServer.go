// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package agentHttpServer

import (
	"fmt"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime/middleware"
	"github.com/kdoctor-io/kdoctor/api/v1/agentServer/server"
	"github.com/kdoctor-io/kdoctor/api/v1/agentServer/server/restapi"
	"github.com/kdoctor-io/kdoctor/api/v1/agentServer/server/restapi/healthy"
	"github.com/kdoctor-io/kdoctor/pkg/types"
	"go.uber.org/zap"
)

// ---------- readiness Healthy Handler
type readinessHealthyHandler struct {
	logger *zap.Logger
}

func (s *readinessHealthyHandler) Handle(r healthy.GetHealthyReadinessParams) middleware.Responder {
	// return healthy.NewGetHealthyReadinessInternalServerError()
	return healthy.NewGetHealthyReadinessOK()
}

// ---------- liveness Healthy Handler
type livenessHealthyHandler struct {
	logger *zap.Logger
}

func (s *livenessHealthyHandler) Handle(r healthy.GetHealthyLivenessParams) middleware.Responder {
	return healthy.NewGetHealthyLivenessOK()
}

// ---------- startup Healthy Handler
type startupHealthyHandler struct {
	logger *zap.Logger
}

func (s *startupHealthyHandler) Handle(r healthy.GetHealthyStartupParams) middleware.Responder {

	return healthy.NewGetHealthyStartupOK()
}

func SetupHealthHttpServer(rootLogger *zap.Logger) {
	logger := rootLogger.Named("health http")

	if types.AgentConfig.AgentHealthPort == 0 {
		logger.Sugar().Warn("health http server is disabled")
		return
	}
	logger.Sugar().Infof("setup agent health http server at port %v", types.AgentConfig.AgentHealthPort)

	spec, err := loads.Embedded(server.SwaggerJSON, server.FlatSwaggerJSON)
	if err != nil {
		logger.Sugar().Fatalf("failed to load Swagger spec, reason=%v ", err)
	}

	api := restapi.NewHTTPServerAPIAPI(spec)
	api.Logger = func(s string, i ...interface{}) {
		logger.Sugar().Infof(s, i)
	}

	// setup route
	api.HealthyGetHealthyReadinessHandler = &readinessHealthyHandler{logger: logger.Named("route: readiness health")}
	api.HealthyGetHealthyLivenessHandler = &livenessHealthyHandler{logger: logger.Named("route: liveness health")}
	api.HealthyGetHealthyStartupHandler = &startupHealthyHandler{logger: logger.Named("route: startup health")}

	srvHealth := server.NewServer(api)
	srvHealth.EnabledListeners = []string{"http"}
	srvHealth.Port = int(types.AgentConfig.AgentHealthPort)
	srvHealth.ConfigureAPI()

	go func() {
		e := srvHealth.Serve()
		s := "health http server break"
		if e != nil {
			s += fmt.Sprintf(" reason=%v", e)
		}
		logger.Fatal(s)
	}()
}
