// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"fmt"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime/middleware"
	"github.com/kdoctor-io/kdoctor/api/v1/agentServer/models"
	"github.com/kdoctor-io/kdoctor/api/v1/agentServer/server"
	"github.com/kdoctor-io/kdoctor/api/v1/agentServer/server/restapi"
	"github.com/kdoctor-io/kdoctor/api/v1/agentServer/server/restapi/echo"
	"github.com/kdoctor-io/kdoctor/api/v1/agentServer/server/restapi/healthy"
	"github.com/kdoctor-io/kdoctor/pkg/types"
	"go.uber.org/zap"
	"os"
)

// ---------- test request Handler

var (
	DefaultInformation = map[string]string{
		"/kdoctoragent": "route to print request",
	}
)

// ---------- "/" route
type echoHandler struct {
	logger *zap.Logger
}

func (s *echoHandler) Handle(r echo.GetParams) middleware.Responder {
	s.logger.Debug("HTTP request from " + r.HTTPRequest.RemoteAddr)

	hostname := types.AgentConfig.PodName
	if len(hostname) == 0 {
		hostname, _ = os.Hostname()
	}
	head := map[string]string{}
	for k, v := range r.HTTPRequest.Header {
		t := ""
		for _, m := range v {
			t += " " + m + " "
		}
		head[k] = t
	}

	t := echo.NewGetOK()
	t.Payload = &models.EchoRes{
		ClientIP:      r.HTTPRequest.RemoteAddr,
		RequestHeader: head,
		RequestURL:    r.HTTPRequest.RequestURI,
		ServerName:    hostname,
		OtherDetail:   DefaultInformation,
	}
	return t
}

// ---------- route "/kdoctoragent"

type echoAgentHandler struct {
	logger *zap.Logger
}

func (s *echoAgentHandler) Handle(r echo.GetParams) middleware.Responder {
	m := r.HTTPRequest
	return (&echoHandler{logger: s.logger}).Handle(echo.GetParams{
		HTTPRequest: m,
	})
}

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

// ====================

func SetupHttpServer() {
	logger := rootLogger.Named("http")

	if types.AgentConfig.HttpPort == 0 {
		logger.Sugar().Warn("http server is disabled")
		return
	}
	logger.Sugar().Infof("setup http server at port %v", types.AgentConfig.HttpPort)

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
	api.EchoGetHandler = &echoHandler{logger: logger.Named("route: request")}
	api.EchoGetHandler = &echoAgentHandler{logger: logger.Named("route: request")}

	//
	srv := server.NewServer(api)
	srv.EnabledListeners = []string{"http"}
	// srv.EnabledListeners = []string{"unix"}
	// srv.SocketPath = "/var/run/http-server-api.sock"

	// default to listen on "0.0.0.0" and "::1"
	// srv.Host = "0.0.0.0"
	srv.Port = int(types.AgentConfig.HttpPort)
	srv.ConfigureAPI()

	go func() {
		e := srv.Serve()
		s := "http server break"
		if e != nil {
			s += fmt.Sprintf(" reason=%v", e)
		}
		logger.Fatal(s)
	}()

}
