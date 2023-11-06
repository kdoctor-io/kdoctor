// Code generated by go-swagger; DO NOT EDIT.

// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package restapi

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/runtime/security"
	"github.com/go-openapi/spec"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"

	"github.com/kdoctor-io/kdoctor/api/v1/agentServer/server/restapi/echo"
	"github.com/kdoctor-io/kdoctor/api/v1/agentServer/server/restapi/healthy"
)

// NewHTTPServerAPIAPI creates a new HTTPServerAPI instance
func NewHTTPServerAPIAPI(spec *loads.Document) *HTTPServerAPIAPI {
	return &HTTPServerAPIAPI{
		handlers:            make(map[string]map[string]http.Handler),
		formats:             strfmt.Default,
		defaultConsumes:     "application/json",
		defaultProduces:     "application/json",
		customConsumers:     make(map[string]runtime.Consumer),
		customProducers:     make(map[string]runtime.Producer),
		PreServerShutdown:   func() {},
		ServerShutdown:      func() {},
		spec:                spec,
		useSwaggerUI:        false,
		ServeError:          errors.ServeError,
		BasicAuthenticator:  security.BasicAuth,
		APIKeyAuthenticator: security.APIKeyAuth,
		BearerAuthenticator: security.BearerAuth,

		JSONConsumer: runtime.JSONConsumer(),

		JSONProducer: runtime.JSONProducer(),

		EchoDeleteHandler: echo.DeleteHandlerFunc(func(params echo.DeleteParams) middleware.Responder {
			return middleware.NotImplemented("operation echo.Delete has not yet been implemented")
		}),
		EchoDeleteKdoctoragentHandler: echo.DeleteKdoctoragentHandlerFunc(func(params echo.DeleteKdoctoragentParams) middleware.Responder {
			return middleware.NotImplemented("operation echo.DeleteKdoctoragent has not yet been implemented")
		}),
		EchoGetHandler: echo.GetHandlerFunc(func(params echo.GetParams) middleware.Responder {
			return middleware.NotImplemented("operation echo.Get has not yet been implemented")
		}),
		HealthyGetHealthyLivenessHandler: healthy.GetHealthyLivenessHandlerFunc(func(params healthy.GetHealthyLivenessParams) middleware.Responder {
			return middleware.NotImplemented("operation healthy.GetHealthyLiveness has not yet been implemented")
		}),
		HealthyGetHealthyReadinessHandler: healthy.GetHealthyReadinessHandlerFunc(func(params healthy.GetHealthyReadinessParams) middleware.Responder {
			return middleware.NotImplemented("operation healthy.GetHealthyReadiness has not yet been implemented")
		}),
		HealthyGetHealthyStartupHandler: healthy.GetHealthyStartupHandlerFunc(func(params healthy.GetHealthyStartupParams) middleware.Responder {
			return middleware.NotImplemented("operation healthy.GetHealthyStartup has not yet been implemented")
		}),
		EchoGetKdoctoragentHandler: echo.GetKdoctoragentHandlerFunc(func(params echo.GetKdoctoragentParams) middleware.Responder {
			return middleware.NotImplemented("operation echo.GetKdoctoragent has not yet been implemented")
		}),
		EchoHeadHandler: echo.HeadHandlerFunc(func(params echo.HeadParams) middleware.Responder {
			return middleware.NotImplemented("operation echo.Head has not yet been implemented")
		}),
		EchoHeadKdoctoragentHandler: echo.HeadKdoctoragentHandlerFunc(func(params echo.HeadKdoctoragentParams) middleware.Responder {
			return middleware.NotImplemented("operation echo.HeadKdoctoragent has not yet been implemented")
		}),
		EchoOptionsHandler: echo.OptionsHandlerFunc(func(params echo.OptionsParams) middleware.Responder {
			return middleware.NotImplemented("operation echo.Options has not yet been implemented")
		}),
		EchoOptionsKdoctoragentHandler: echo.OptionsKdoctoragentHandlerFunc(func(params echo.OptionsKdoctoragentParams) middleware.Responder {
			return middleware.NotImplemented("operation echo.OptionsKdoctoragent has not yet been implemented")
		}),
		EchoPatchHandler: echo.PatchHandlerFunc(func(params echo.PatchParams) middleware.Responder {
			return middleware.NotImplemented("operation echo.Patch has not yet been implemented")
		}),
		EchoPatchKdoctoragentHandler: echo.PatchKdoctoragentHandlerFunc(func(params echo.PatchKdoctoragentParams) middleware.Responder {
			return middleware.NotImplemented("operation echo.PatchKdoctoragent has not yet been implemented")
		}),
		EchoPostHandler: echo.PostHandlerFunc(func(params echo.PostParams) middleware.Responder {
			return middleware.NotImplemented("operation echo.Post has not yet been implemented")
		}),
		EchoPostKdoctoragentHandler: echo.PostKdoctoragentHandlerFunc(func(params echo.PostKdoctoragentParams) middleware.Responder {
			return middleware.NotImplemented("operation echo.PostKdoctoragent has not yet been implemented")
		}),
		EchoPutHandler: echo.PutHandlerFunc(func(params echo.PutParams) middleware.Responder {
			return middleware.NotImplemented("operation echo.Put has not yet been implemented")
		}),
		EchoPutKdoctoragentHandler: echo.PutKdoctoragentHandlerFunc(func(params echo.PutKdoctoragentParams) middleware.Responder {
			return middleware.NotImplemented("operation echo.PutKdoctoragent has not yet been implemented")
		}),
	}
}

/*HTTPServerAPIAPI agent http server */
type HTTPServerAPIAPI struct {
	spec            *loads.Document
	context         *middleware.Context
	handlers        map[string]map[string]http.Handler
	formats         strfmt.Registry
	customConsumers map[string]runtime.Consumer
	customProducers map[string]runtime.Producer
	defaultConsumes string
	defaultProduces string
	Middleware      func(middleware.Builder) http.Handler
	useSwaggerUI    bool

	// BasicAuthenticator generates a runtime.Authenticator from the supplied basic auth function.
	// It has a default implementation in the security package, however you can replace it for your particular usage.
	BasicAuthenticator func(security.UserPassAuthentication) runtime.Authenticator

	// APIKeyAuthenticator generates a runtime.Authenticator from the supplied token auth function.
	// It has a default implementation in the security package, however you can replace it for your particular usage.
	APIKeyAuthenticator func(string, string, security.TokenAuthentication) runtime.Authenticator

	// BearerAuthenticator generates a runtime.Authenticator from the supplied bearer token auth function.
	// It has a default implementation in the security package, however you can replace it for your particular usage.
	BearerAuthenticator func(string, security.ScopedTokenAuthentication) runtime.Authenticator

	// JSONConsumer registers a consumer for the following mime types:
	//   - application/json
	JSONConsumer runtime.Consumer

	// JSONProducer registers a producer for the following mime types:
	//   - application/json
	JSONProducer runtime.Producer

	// EchoDeleteHandler sets the operation handler for the delete operation
	EchoDeleteHandler echo.DeleteHandler
	// EchoDeleteKdoctoragentHandler sets the operation handler for the delete kdoctoragent operation
	EchoDeleteKdoctoragentHandler echo.DeleteKdoctoragentHandler
	// EchoGetHandler sets the operation handler for the get operation
	EchoGetHandler echo.GetHandler
	// HealthyGetHealthyLivenessHandler sets the operation handler for the get healthy liveness operation
	HealthyGetHealthyLivenessHandler healthy.GetHealthyLivenessHandler
	// HealthyGetHealthyReadinessHandler sets the operation handler for the get healthy readiness operation
	HealthyGetHealthyReadinessHandler healthy.GetHealthyReadinessHandler
	// HealthyGetHealthyStartupHandler sets the operation handler for the get healthy startup operation
	HealthyGetHealthyStartupHandler healthy.GetHealthyStartupHandler
	// EchoGetKdoctoragentHandler sets the operation handler for the get kdoctoragent operation
	EchoGetKdoctoragentHandler echo.GetKdoctoragentHandler
	// EchoHeadHandler sets the operation handler for the head operation
	EchoHeadHandler echo.HeadHandler
	// EchoHeadKdoctoragentHandler sets the operation handler for the head kdoctoragent operation
	EchoHeadKdoctoragentHandler echo.HeadKdoctoragentHandler
	// EchoOptionsHandler sets the operation handler for the options operation
	EchoOptionsHandler echo.OptionsHandler
	// EchoOptionsKdoctoragentHandler sets the operation handler for the options kdoctoragent operation
	EchoOptionsKdoctoragentHandler echo.OptionsKdoctoragentHandler
	// EchoPatchHandler sets the operation handler for the patch operation
	EchoPatchHandler echo.PatchHandler
	// EchoPatchKdoctoragentHandler sets the operation handler for the patch kdoctoragent operation
	EchoPatchKdoctoragentHandler echo.PatchKdoctoragentHandler
	// EchoPostHandler sets the operation handler for the post operation
	EchoPostHandler echo.PostHandler
	// EchoPostKdoctoragentHandler sets the operation handler for the post kdoctoragent operation
	EchoPostKdoctoragentHandler echo.PostKdoctoragentHandler
	// EchoPutHandler sets the operation handler for the put operation
	EchoPutHandler echo.PutHandler
	// EchoPutKdoctoragentHandler sets the operation handler for the put kdoctoragent operation
	EchoPutKdoctoragentHandler echo.PutKdoctoragentHandler

	// ServeError is called when an error is received, there is a default handler
	// but you can set your own with this
	ServeError func(http.ResponseWriter, *http.Request, error)

	// PreServerShutdown is called before the HTTP(S) server is shutdown
	// This allows for custom functions to get executed before the HTTP(S) server stops accepting traffic
	PreServerShutdown func()

	// ServerShutdown is called when the HTTP(S) server is shut down and done
	// handling all active connections and does not accept connections any more
	ServerShutdown func()

	// Custom command line argument groups with their descriptions
	CommandLineOptionsGroups []swag.CommandLineOptionsGroup

	// User defined logger function.
	Logger func(string, ...interface{})
}

// UseRedoc for documentation at /docs
func (o *HTTPServerAPIAPI) UseRedoc() {
	o.useSwaggerUI = false
}

// UseSwaggerUI for documentation at /docs
func (o *HTTPServerAPIAPI) UseSwaggerUI() {
	o.useSwaggerUI = true
}

// SetDefaultProduces sets the default produces media type
func (o *HTTPServerAPIAPI) SetDefaultProduces(mediaType string) {
	o.defaultProduces = mediaType
}

// SetDefaultConsumes returns the default consumes media type
func (o *HTTPServerAPIAPI) SetDefaultConsumes(mediaType string) {
	o.defaultConsumes = mediaType
}

// SetSpec sets a spec that will be served for the clients.
func (o *HTTPServerAPIAPI) SetSpec(spec *loads.Document) {
	o.spec = spec
}

// DefaultProduces returns the default produces media type
func (o *HTTPServerAPIAPI) DefaultProduces() string {
	return o.defaultProduces
}

// DefaultConsumes returns the default consumes media type
func (o *HTTPServerAPIAPI) DefaultConsumes() string {
	return o.defaultConsumes
}

// Formats returns the registered string formats
func (o *HTTPServerAPIAPI) Formats() strfmt.Registry {
	return o.formats
}

// RegisterFormat registers a custom format validator
func (o *HTTPServerAPIAPI) RegisterFormat(name string, format strfmt.Format, validator strfmt.Validator) {
	o.formats.Add(name, format, validator)
}

// Validate validates the registrations in the HTTPServerAPIAPI
func (o *HTTPServerAPIAPI) Validate() error {
	var unregistered []string

	if o.JSONConsumer == nil {
		unregistered = append(unregistered, "JSONConsumer")
	}

	if o.JSONProducer == nil {
		unregistered = append(unregistered, "JSONProducer")
	}

	if o.EchoDeleteHandler == nil {
		unregistered = append(unregistered, "echo.DeleteHandler")
	}
	if o.EchoDeleteKdoctoragentHandler == nil {
		unregistered = append(unregistered, "echo.DeleteKdoctoragentHandler")
	}
	if o.EchoGetHandler == nil {
		unregistered = append(unregistered, "echo.GetHandler")
	}
	if o.HealthyGetHealthyLivenessHandler == nil {
		unregistered = append(unregistered, "healthy.GetHealthyLivenessHandler")
	}
	if o.HealthyGetHealthyReadinessHandler == nil {
		unregistered = append(unregistered, "healthy.GetHealthyReadinessHandler")
	}
	if o.HealthyGetHealthyStartupHandler == nil {
		unregistered = append(unregistered, "healthy.GetHealthyStartupHandler")
	}
	if o.EchoGetKdoctoragentHandler == nil {
		unregistered = append(unregistered, "echo.GetKdoctoragentHandler")
	}
	if o.EchoHeadHandler == nil {
		unregistered = append(unregistered, "echo.HeadHandler")
	}
	if o.EchoHeadKdoctoragentHandler == nil {
		unregistered = append(unregistered, "echo.HeadKdoctoragentHandler")
	}
	if o.EchoOptionsHandler == nil {
		unregistered = append(unregistered, "echo.OptionsHandler")
	}
	if o.EchoOptionsKdoctoragentHandler == nil {
		unregistered = append(unregistered, "echo.OptionsKdoctoragentHandler")
	}
	if o.EchoPatchHandler == nil {
		unregistered = append(unregistered, "echo.PatchHandler")
	}
	if o.EchoPatchKdoctoragentHandler == nil {
		unregistered = append(unregistered, "echo.PatchKdoctoragentHandler")
	}
	if o.EchoPostHandler == nil {
		unregistered = append(unregistered, "echo.PostHandler")
	}
	if o.EchoPostKdoctoragentHandler == nil {
		unregistered = append(unregistered, "echo.PostKdoctoragentHandler")
	}
	if o.EchoPutHandler == nil {
		unregistered = append(unregistered, "echo.PutHandler")
	}
	if o.EchoPutKdoctoragentHandler == nil {
		unregistered = append(unregistered, "echo.PutKdoctoragentHandler")
	}

	if len(unregistered) > 0 {
		return fmt.Errorf("missing registration: %s", strings.Join(unregistered, ", "))
	}

	return nil
}

// ServeErrorFor gets a error handler for a given operation id
func (o *HTTPServerAPIAPI) ServeErrorFor(operationID string) func(http.ResponseWriter, *http.Request, error) {
	return o.ServeError
}

// AuthenticatorsFor gets the authenticators for the specified security schemes
func (o *HTTPServerAPIAPI) AuthenticatorsFor(schemes map[string]spec.SecurityScheme) map[string]runtime.Authenticator {
	return nil
}

// Authorizer returns the registered authorizer
func (o *HTTPServerAPIAPI) Authorizer() runtime.Authorizer {
	return nil
}

// ConsumersFor gets the consumers for the specified media types.
// MIME type parameters are ignored here.
func (o *HTTPServerAPIAPI) ConsumersFor(mediaTypes []string) map[string]runtime.Consumer {
	result := make(map[string]runtime.Consumer, len(mediaTypes))
	for _, mt := range mediaTypes {
		switch mt {
		case "application/json":
			result["application/json"] = o.JSONConsumer
		}

		if c, ok := o.customConsumers[mt]; ok {
			result[mt] = c
		}
	}
	return result
}

// ProducersFor gets the producers for the specified media types.
// MIME type parameters are ignored here.
func (o *HTTPServerAPIAPI) ProducersFor(mediaTypes []string) map[string]runtime.Producer {
	result := make(map[string]runtime.Producer, len(mediaTypes))
	for _, mt := range mediaTypes {
		switch mt {
		case "application/json":
			result["application/json"] = o.JSONProducer
		}

		if p, ok := o.customProducers[mt]; ok {
			result[mt] = p
		}
	}
	return result
}

// HandlerFor gets a http.Handler for the provided operation method and path
func (o *HTTPServerAPIAPI) HandlerFor(method, path string) (http.Handler, bool) {
	if o.handlers == nil {
		return nil, false
	}
	um := strings.ToUpper(method)
	if _, ok := o.handlers[um]; !ok {
		return nil, false
	}
	if path == "/" {
		path = ""
	}
	h, ok := o.handlers[um][path]
	return h, ok
}

// Context returns the middleware context for the HTTP server API API
func (o *HTTPServerAPIAPI) Context() *middleware.Context {
	if o.context == nil {
		o.context = middleware.NewRoutableContext(o.spec, o, nil)
	}

	return o.context
}

func (o *HTTPServerAPIAPI) initHandlerCache() {
	o.Context() // don't care about the result, just that the initialization happened
	if o.handlers == nil {
		o.handlers = make(map[string]map[string]http.Handler)
	}

	if o.handlers["DELETE"] == nil {
		o.handlers["DELETE"] = make(map[string]http.Handler)
	}
	o.handlers["DELETE"][""] = echo.NewDelete(o.context, o.EchoDeleteHandler)
	if o.handlers["DELETE"] == nil {
		o.handlers["DELETE"] = make(map[string]http.Handler)
	}
	o.handlers["DELETE"]["/kdoctoragent"] = echo.NewDeleteKdoctoragent(o.context, o.EchoDeleteKdoctoragentHandler)
	if o.handlers["GET"] == nil {
		o.handlers["GET"] = make(map[string]http.Handler)
	}
	o.handlers["GET"][""] = echo.NewGet(o.context, o.EchoGetHandler)
	if o.handlers["GET"] == nil {
		o.handlers["GET"] = make(map[string]http.Handler)
	}
	o.handlers["GET"]["/healthy/liveness"] = healthy.NewGetHealthyLiveness(o.context, o.HealthyGetHealthyLivenessHandler)
	if o.handlers["GET"] == nil {
		o.handlers["GET"] = make(map[string]http.Handler)
	}
	o.handlers["GET"]["/healthy/readiness"] = healthy.NewGetHealthyReadiness(o.context, o.HealthyGetHealthyReadinessHandler)
	if o.handlers["GET"] == nil {
		o.handlers["GET"] = make(map[string]http.Handler)
	}
	o.handlers["GET"]["/healthy/startup"] = healthy.NewGetHealthyStartup(o.context, o.HealthyGetHealthyStartupHandler)
	if o.handlers["GET"] == nil {
		o.handlers["GET"] = make(map[string]http.Handler)
	}
	o.handlers["GET"]["/kdoctoragent"] = echo.NewGetKdoctoragent(o.context, o.EchoGetKdoctoragentHandler)
	if o.handlers["HEAD"] == nil {
		o.handlers["HEAD"] = make(map[string]http.Handler)
	}
	o.handlers["HEAD"][""] = echo.NewHead(o.context, o.EchoHeadHandler)
	if o.handlers["HEAD"] == nil {
		o.handlers["HEAD"] = make(map[string]http.Handler)
	}
	o.handlers["HEAD"]["/kdoctoragent"] = echo.NewHeadKdoctoragent(o.context, o.EchoHeadKdoctoragentHandler)
	if o.handlers["OPTIONS"] == nil {
		o.handlers["OPTIONS"] = make(map[string]http.Handler)
	}
	o.handlers["OPTIONS"][""] = echo.NewOptions(o.context, o.EchoOptionsHandler)
	if o.handlers["OPTIONS"] == nil {
		o.handlers["OPTIONS"] = make(map[string]http.Handler)
	}
	o.handlers["OPTIONS"]["/kdoctoragent"] = echo.NewOptionsKdoctoragent(o.context, o.EchoOptionsKdoctoragentHandler)
	if o.handlers["PATCH"] == nil {
		o.handlers["PATCH"] = make(map[string]http.Handler)
	}
	o.handlers["PATCH"][""] = echo.NewPatch(o.context, o.EchoPatchHandler)
	if o.handlers["PATCH"] == nil {
		o.handlers["PATCH"] = make(map[string]http.Handler)
	}
	o.handlers["PATCH"]["/kdoctoragent"] = echo.NewPatchKdoctoragent(o.context, o.EchoPatchKdoctoragentHandler)
	if o.handlers["POST"] == nil {
		o.handlers["POST"] = make(map[string]http.Handler)
	}
	o.handlers["POST"][""] = echo.NewPost(o.context, o.EchoPostHandler)
	if o.handlers["POST"] == nil {
		o.handlers["POST"] = make(map[string]http.Handler)
	}
	o.handlers["POST"]["/kdoctoragent"] = echo.NewPostKdoctoragent(o.context, o.EchoPostKdoctoragentHandler)
	if o.handlers["PUT"] == nil {
		o.handlers["PUT"] = make(map[string]http.Handler)
	}
	o.handlers["PUT"][""] = echo.NewPut(o.context, o.EchoPutHandler)
	if o.handlers["PUT"] == nil {
		o.handlers["PUT"] = make(map[string]http.Handler)
	}
	o.handlers["PUT"]["/kdoctoragent"] = echo.NewPutKdoctoragent(o.context, o.EchoPutKdoctoragentHandler)
}

// Serve creates a http handler to serve the API over HTTP
// can be used directly in http.ListenAndServe(":8000", api.Serve(nil))
func (o *HTTPServerAPIAPI) Serve(builder middleware.Builder) http.Handler {
	o.Init()

	if o.Middleware != nil {
		return o.Middleware(builder)
	}
	if o.useSwaggerUI {
		return o.context.APIHandlerSwaggerUI(builder)
	}
	return o.context.APIHandler(builder)
}

// Init allows you to just initialize the handler cache, you can then recompose the middleware as you see fit
func (o *HTTPServerAPIAPI) Init() {
	if len(o.handlers) == 0 {
		o.initHandlerCache()
	}
}

// RegisterConsumer allows you to add (or override) a consumer for a media type.
func (o *HTTPServerAPIAPI) RegisterConsumer(mediaType string, consumer runtime.Consumer) {
	o.customConsumers[mediaType] = consumer
}

// RegisterProducer allows you to add (or override) a producer for a media type.
func (o *HTTPServerAPIAPI) RegisterProducer(mediaType string, producer runtime.Producer) {
	o.customProducers[mediaType] = producer
}

// AddMiddlewareFor adds a http middleware to existing handler
func (o *HTTPServerAPIAPI) AddMiddlewareFor(method, path string, builder middleware.Builder) {
	um := strings.ToUpper(method)
	if path == "/" {
		path = ""
	}
	o.Init()
	if h, ok := o.handlers[um][path]; ok {
		o.handlers[method][path] = builder(h)
	}
}