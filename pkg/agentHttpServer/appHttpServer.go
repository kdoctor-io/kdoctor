// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package agentHttpServer

import (
	"fmt"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime/middleware"
	"github.com/jessevdk/go-flags"
	"github.com/kdoctor-io/kdoctor/api/v1/agentServer/models"
	"github.com/kdoctor-io/kdoctor/api/v1/agentServer/server"
	"github.com/kdoctor-io/kdoctor/api/v1/agentServer/server/restapi"
	"github.com/kdoctor-io/kdoctor/api/v1/agentServer/server/restapi/echo"
	"github.com/kdoctor-io/kdoctor/pkg/lock"
	"github.com/kdoctor-io/kdoctor/pkg/types"
	"go.uber.org/zap"
	"os"
	"time"
)

var (
	ParamInformation = map[string]string{
		"delay": "in query, delay some second return response",
	}
	SupportedMethod = []string{"GET", "PUT", "POST", "DELETE", "HEAD", "PATCH", "OPTIONS"}

	RequestCounts = newRequestCount()
)

type requestCounts struct {
	counts map[string]int64
	l      lock.Mutex
}

func newRequestCount() *requestCounts {
	return &requestCounts{
		counts: make(map[string]int64),
		l:      lock.Mutex{},
	}
}
func (rc *requestCounts) AddOneCount(task string) {
	rc.l.Lock()
	rc.counts[task]++
	rc.l.Unlock()
}

func (rc *requestCounts) ZeroCount(task string) {
	rc.l.Lock()
	rc.counts[task] = 0
	rc.l.Unlock()
}

func (rc *requestCounts) GetCount(task string) int64 {
	rc.l.Lock()
	r := rc.counts[task]
	rc.l.Unlock()
	return r
}

// route /
// ---------- get
type echoGetHandler struct {
	logger *zap.Logger
}

func (s *echoGetHandler) Handle(r echo.GetParams) middleware.Responder {
	if r.Delay != nil {
		s.logger.Sugar().Debugf("%s method  %s delay %d request from %s", r.HTTPRequest.Proto, r.HTTPRequest.Method, *r.Delay, r.HTTPRequest.RemoteAddr)
	} else {
		s.logger.Sugar().Debugf("%s method  %s delay 0 request from %s", r.HTTPRequest.Proto, r.HTTPRequest.Method, r.HTTPRequest.RemoteAddr)
	}

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
	var task string
	if r.Task == nil {
		task = "default"
	} else {
		task = *r.Task
	}
	RequestCounts.AddOneCount(task)

	t := echo.NewGetOK()
	t.Payload = &models.EchoRes{
		ClientIP:        r.HTTPRequest.RemoteAddr,
		RequestHeader:   head,
		RequestURL:      r.HTTPRequest.RequestURI,
		RequestMethod:   r.HTTPRequest.Method,
		ServerName:      hostname,
		RequestCount:    RequestCounts.GetCount(task),
		ParamDetail:     ParamInformation,
		SupportedMethod: SupportedMethod,
		TaskName:        task,
	}
	if r.Delay != nil {
		time.Sleep(time.Duration(*r.Delay) * time.Second)
		t.Payload.RequestParam = fmt.Sprintf("delay=%d", *r.Delay)
	}
	return t
}

// -----------delete
type echoDeleteHandler struct {
	logger *zap.Logger
}

func (s *echoDeleteHandler) Handle(r echo.DeleteParams) middleware.Responder {
	var task string
	if r.Task == nil {
		task = "default"
	} else {
		task = *r.Task
	}
	RequestCounts.ZeroCount(task)
	return echo.NewDeleteOK()
}

// ----------- put
type echoPutHandler struct {
	logger *zap.Logger
}

func (s *echoPutHandler) Handle(r echo.PutParams) middleware.Responder {
	if r.Delay != nil {
		s.logger.Sugar().Debugf("%s method  %s delay %d request from %s", r.HTTPRequest.Proto, r.HTTPRequest.Method, *r.Delay, r.HTTPRequest.RemoteAddr)
	} else {
		s.logger.Sugar().Debugf("%s method  %s delay 0 request from %s", r.HTTPRequest.Proto, r.HTTPRequest.Method, r.HTTPRequest.RemoteAddr)
	}
	var task string
	if r.Task == nil {
		task = "default"
	} else {
		task = *r.Task
	}
	RequestCounts.AddOneCount(task)
	if r.Delay != nil {
		time.Sleep(time.Duration(*r.Delay) * time.Second)
	}
	return echo.NewPutOK()
}

// ----------- post
type echoPostHandler struct {
	logger *zap.Logger
}

func (s *echoPostHandler) Handle(r echo.PostParams) middleware.Responder {
	if r.Delay != nil {
		s.logger.Sugar().Debugf("%s method  %s delay %d request from %s", r.HTTPRequest.Proto, r.HTTPRequest.Method, *r.Delay, r.HTTPRequest.RemoteAddr)
	} else {
		s.logger.Sugar().Debugf("%s method  %s delay 0 request from %s", r.HTTPRequest.Proto, r.HTTPRequest.Method, r.HTTPRequest.RemoteAddr)
	}
	var task string
	if r.Task == nil {
		task = "default"
	} else {
		task = *r.Task
	}
	RequestCounts.AddOneCount(task)
	if r.Delay != nil {
		time.Sleep(time.Duration(*r.Delay) * time.Second)
	}
	_, err := r.TestArgs.MarshalBinary()
	if err != nil {
		return echo.NewPostInternalServerError()
	}
	return echo.NewPostOK()
}

// ----------- head
type echoHeadHandler struct {
	logger *zap.Logger
}

func (s *echoHeadHandler) Handle(r echo.HeadParams) middleware.Responder {
	if r.Delay != nil {
		s.logger.Sugar().Debugf("%s method  %s delay %d request from %s", r.HTTPRequest.Proto, r.HTTPRequest.Method, *r.Delay, r.HTTPRequest.RemoteAddr)
	} else {
		s.logger.Sugar().Debugf("%s method  %s delay 0 request from %s", r.HTTPRequest.Proto, r.HTTPRequest.Method, r.HTTPRequest.RemoteAddr)
	}
	var task string
	if r.Task == nil {
		task = "default"
	} else {
		task = *r.Task
	}
	RequestCounts.AddOneCount(task)
	if r.Delay != nil {
		time.Sleep(time.Duration(*r.Delay) * time.Second)
	}
	return echo.NewHeadOK()
}

// ----------- options
type echoOptionsHandler struct {
	logger *zap.Logger
}

func (s *echoOptionsHandler) Handle(r echo.OptionsParams) middleware.Responder {
	if r.Delay != nil {
		s.logger.Sugar().Debugf("%s method  %s delay %d request from %s", r.HTTPRequest.Proto, r.HTTPRequest.Method, *r.Delay, r.HTTPRequest.RemoteAddr)
	} else {
		s.logger.Sugar().Debugf("%s method  %s delay 0 request from %s", r.HTTPRequest.Proto, r.HTTPRequest.Method, r.HTTPRequest.RemoteAddr)
	}
	var task string
	if r.Task == nil {
		task = "default"
	} else {
		task = *r.Task
	}
	RequestCounts.AddOneCount(task)
	if r.Delay != nil {
		time.Sleep(time.Duration(*r.Delay) * time.Second)
	}
	return echo.NewOptionsOK()
}

// ----------- patch
type echoPatchHandler struct {
	logger *zap.Logger
}

func (s *echoPatchHandler) Handle(r echo.PatchParams) middleware.Responder {
	if r.Delay != nil {
		s.logger.Sugar().Debugf("%s method  %s delay %d request from %s", r.HTTPRequest.Proto, r.HTTPRequest.Method, *r.Delay, r.HTTPRequest.RemoteAddr)
	} else {
		s.logger.Sugar().Debugf("%s method  %s delay 0 request from %s", r.HTTPRequest.Proto, r.HTTPRequest.Method, r.HTTPRequest.RemoteAddr)
	}
	var task string
	if r.Task == nil {
		task = "default"
	} else {
		task = *r.Task
	}
	RequestCounts.AddOneCount(task)
	if r.Delay != nil {
		time.Sleep(time.Duration(*r.Delay) * time.Second)
	}
	return echo.NewPatchOK()
}

// route /kdoctoragent
// ---------- get
type echoKdoctorGetHandler struct {
	logger *zap.Logger
}

func (s *echoKdoctorGetHandler) Handle(r echo.GetKdoctoragentParams) middleware.Responder {
	if r.Delay != nil {
		s.logger.Sugar().Debugf("%s method  %s delay %d request from %s", r.HTTPRequest.Proto, r.HTTPRequest.Method, *r.Delay, r.HTTPRequest.RemoteAddr)
	} else {
		s.logger.Sugar().Debugf("%s method  %s delay 0 request from %s", r.HTTPRequest.Proto, r.HTTPRequest.Method, r.HTTPRequest.RemoteAddr)
	}

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
	var task string
	if r.Task == nil {
		task = "default"
	} else {
		task = *r.Task
	}
	RequestCounts.AddOneCount(task)
	t := echo.NewGetKdoctoragentOK()
	t.Payload = &models.EchoRes{
		ClientIP:        r.HTTPRequest.RemoteAddr,
		RequestHeader:   head,
		RequestURL:      r.HTTPRequest.RequestURI,
		RequestMethod:   r.HTTPRequest.Method,
		ServerName:      hostname,
		RequestCount:    RequestCounts.GetCount(task),
		ParamDetail:     ParamInformation,
		SupportedMethod: SupportedMethod,
		TaskName:        task,
	}
	if r.Delay != nil {
		time.Sleep(time.Duration(*r.Delay) * time.Second)
		t.Payload.RequestParam = fmt.Sprintf("delay=%d", *r.Delay)
	}
	return t
}

// -----------delete
type echoKdoctorDeleteHandler struct {
	logger *zap.Logger
}

func (s *echoKdoctorDeleteHandler) Handle(r echo.DeleteKdoctoragentParams) middleware.Responder {
	var task string
	if r.Task == nil {
		task = "default"
	} else {
		task = *r.Task
	}
	RequestCounts.ZeroCount(task)
	return echo.NewDeleteKdoctoragentOK()
}

// ----------- put
type echoKdoctorPutHandler struct {
	logger *zap.Logger
}

func (s *echoKdoctorPutHandler) Handle(r echo.PutKdoctoragentParams) middleware.Responder {
	if r.Delay != nil {
		s.logger.Sugar().Debugf("%s method  %s delay %d request from %s", r.HTTPRequest.Proto, r.HTTPRequest.Method, *r.Delay, r.HTTPRequest.RemoteAddr)
	} else {
		s.logger.Sugar().Debugf("%s method  %s delay 0 request from %s", r.HTTPRequest.Proto, r.HTTPRequest.Method, r.HTTPRequest.RemoteAddr)
	}
	var task string
	if r.Task == nil {
		task = "default"
	} else {
		task = *r.Task
	}
	RequestCounts.AddOneCount(task)
	if r.Delay != nil {
		time.Sleep(time.Duration(*r.Delay) * time.Second)
	}
	return echo.NewPutKdoctoragentOK()
}

// ----------- post
type echoKdoctorPostHandler struct {
	logger *zap.Logger
}

func (s *echoKdoctorPostHandler) Handle(r echo.PostKdoctoragentParams) middleware.Responder {
	if r.Delay != nil {
		s.logger.Sugar().Debugf("%s method  %s delay %d request from %s", r.HTTPRequest.Proto, r.HTTPRequest.Method, *r.Delay, r.HTTPRequest.RemoteAddr)
	} else {
		s.logger.Sugar().Debugf("%s method  %s delay 0 request from %s", r.HTTPRequest.Proto, r.HTTPRequest.Method, r.HTTPRequest.RemoteAddr)
	}
	var task string
	if r.Task == nil {
		task = "default"
	} else {
		task = *r.Task
	}
	RequestCounts.AddOneCount(task)
	if r.Delay != nil {
		time.Sleep(time.Duration(*r.Delay) * time.Second)
	}
	_, err := r.TestArgs.MarshalBinary()
	if err != nil {
		return echo.NewPostKdoctoragentInternalServerError()
	}
	return echo.NewPostKdoctoragentOK()
}

// ----------- head
type echoKdoctorHeadHandler struct {
	logger *zap.Logger
}

func (s *echoKdoctorHeadHandler) Handle(r echo.HeadKdoctoragentParams) middleware.Responder {
	if r.Delay != nil {
		s.logger.Sugar().Debugf("%s method  %s delay %d request from %s", r.HTTPRequest.Proto, r.HTTPRequest.Method, *r.Delay, r.HTTPRequest.RemoteAddr)
	} else {
		s.logger.Sugar().Debugf("%s method  %s delay 0 request from %s", r.HTTPRequest.Proto, r.HTTPRequest.Method, r.HTTPRequest.RemoteAddr)
	}
	var task string
	if r.Task == nil {
		task = "default"
	} else {
		task = *r.Task
	}
	RequestCounts.AddOneCount(task)
	if r.Delay != nil {
		time.Sleep(time.Duration(*r.Delay) * time.Second)
	}
	return echo.NewHeadKdoctoragentOK()
}

// ----------- options
type echoKdoctorOptionsHandler struct {
	logger *zap.Logger
}

func (s *echoKdoctorOptionsHandler) Handle(r echo.OptionsKdoctoragentParams) middleware.Responder {
	if r.Delay != nil {
		s.logger.Sugar().Debugf("%s method  %s delay %d request from %s", r.HTTPRequest.Proto, r.HTTPRequest.Method, *r.Delay, r.HTTPRequest.RemoteAddr)
	} else {
		s.logger.Sugar().Debugf("%s method  %s delay 0 request from %s", r.HTTPRequest.Proto, r.HTTPRequest.Method, r.HTTPRequest.RemoteAddr)
	}
	var task string
	if r.Task == nil {
		task = "default"
	} else {
		task = *r.Task
	}
	RequestCounts.AddOneCount(task)
	if r.Delay != nil {
		time.Sleep(time.Duration(*r.Delay) * time.Second)
	}
	return echo.NewOptionsKdoctoragentOK()
}

// ----------- patch
type echoKdoctorPatchHandler struct {
	logger *zap.Logger
}

func (s *echoKdoctorPatchHandler) Handle(r echo.PatchKdoctoragentParams) middleware.Responder {
	if r.Delay != nil {
		s.logger.Sugar().Debugf("%s method  %s delay %d request from %s", r.HTTPRequest.Proto, r.HTTPRequest.Method, *r.Delay, r.HTTPRequest.RemoteAddr)
	} else {
		s.logger.Sugar().Debugf("%s method  %s delay 0 request from %s", r.HTTPRequest.Proto, r.HTTPRequest.Method, r.HTTPRequest.RemoteAddr)
	}
	var task string
	if r.Task == nil {
		task = "default"
	} else {
		task = *r.Task
	}
	RequestCounts.AddOneCount(task)
	if r.Delay != nil {
		time.Sleep(time.Duration(*r.Delay) * time.Second)
	}
	return echo.NewPatchKdoctoragentOK()
}

func SetupAppHttpServer(rootLogger *zap.Logger, tlsCert, tlsKey string) {
	logger := rootLogger.Named("app http")

	if types.AgentConfig.AppHttpPort == 0 {
		logger.Sugar().Warn("app http server is disabled")
		return
	}

	spec, err := loads.Embedded(server.SwaggerJSON, server.FlatSwaggerJSON)
	if err != nil {
		logger.Sugar().Fatalf("failed to load Swagger spec, reason=%v ", err)
	}

	api := restapi.NewHTTPServerAPIAPI(spec)
	api.Logger = func(s string, i ...interface{}) {
		logger.Sugar().Infof(s, i)
	}

	// setup route "/"
	api.EchoGetHandler = &echoGetHandler{logger: logger.Named("route: request")}
	api.EchoDeleteHandler = &echoDeleteHandler{logger: logger.Named("route: summary counts")}
	api.EchoPostHandler = &echoPostHandler{logger: logger.Named("route: summary counts")}
	api.EchoPutHandler = &echoPutHandler{logger: logger.Named("route: summary counts")}
	api.EchoHeadHandler = &echoHeadHandler{logger: logger.Named("route: summary counts")}
	api.EchoOptionsHandler = &echoOptionsHandler{logger: logger.Named("route: summary counts")}
	api.EchoPatchHandler = &echoPatchHandler{logger: logger.Named("route: summary counts")}

	// setup route "/kdoctoragent"
	api.EchoGetKdoctoragentHandler = &echoKdoctorGetHandler{logger: logger.Named("route: request")}
	api.EchoDeleteKdoctoragentHandler = &echoKdoctorDeleteHandler{logger: logger.Named("route: summary counts")}
	api.EchoPostKdoctoragentHandler = &echoKdoctorPostHandler{logger: logger.Named("route: summary counts")}
	api.EchoPutKdoctoragentHandler = &echoKdoctorPutHandler{logger: logger.Named("route: summary counts")}
	api.EchoHeadKdoctoragentHandler = &echoKdoctorHeadHandler{logger: logger.Named("route: summary counts")}
	api.EchoOptionsKdoctoragentHandler = &echoKdoctorOptionsHandler{logger: logger.Named("route: summary counts")}
	api.EchoPatchKdoctoragentHandler = &echoKdoctorPatchHandler{logger: logger.Named("route: summary counts")}

	srvApp := server.NewServer(api)
	srvApp.EnabledListeners = []string{"https", "http"}
	// http
	srvApp.Port = int(types.AgentConfig.AppHttpPort)
	// https
	srvApp.TLSPort = int(types.AgentConfig.AppHttpsPort)
	// verify ca
	if !types.AgentConfig.TlsInsecure {
		srvApp.TLSCACertificate = flags.Filename(types.AgentConfig.TlsCaCertPath)
		logger.Sugar().Infof("agent enabled verify tls")
	} else {
		logger.Sugar().Infof("agent disabled verify tls")
	}
	srvApp.TLSCertificate = flags.Filename(tlsCert)
	srvApp.TLSCertificateKey = flags.Filename(tlsKey)

	logger.Sugar().Infof("setup agent app http server at port %v", types.AgentConfig.AppHttpPort)
	logger.Sugar().Infof("setup agent app https server at port %v", types.AgentConfig.AppHttpsPort)

	srvApp.ConfigureAPI()
	go func() {
		e := srvApp.Serve()
		s := "app http server break"
		if e != nil {
			s += fmt.Sprintf(" reason=%v", e)
		}
		logger.Fatal(s)
	}()
}
