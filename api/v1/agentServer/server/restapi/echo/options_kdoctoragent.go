// Code generated by go-swagger; DO NOT EDIT.

// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package echo

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// OptionsKdoctoragentHandlerFunc turns a function with the right signature into a options kdoctoragent handler
type OptionsKdoctoragentHandlerFunc func(OptionsKdoctoragentParams) middleware.Responder

// Handle executing the request and returning a response
func (fn OptionsKdoctoragentHandlerFunc) Handle(params OptionsKdoctoragentParams) middleware.Responder {
	return fn(params)
}

// OptionsKdoctoragentHandler interface for that can handle valid options kdoctoragent params
type OptionsKdoctoragentHandler interface {
	Handle(OptionsKdoctoragentParams) middleware.Responder
}

// NewOptionsKdoctoragent creates a new http.Handler for the options kdoctoragent operation
func NewOptionsKdoctoragent(ctx *middleware.Context, handler OptionsKdoctoragentHandler) *OptionsKdoctoragent {
	return &OptionsKdoctoragent{Context: ctx, Handler: handler}
}

/*
	OptionsKdoctoragent swagger:route OPTIONS /kdoctoragent echo optionsKdoctoragent

echo http request

echo http request
*/
type OptionsKdoctoragent struct {
	Context *middleware.Context
	Handler OptionsKdoctoragentHandler
}

func (o *OptionsKdoctoragent) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewOptionsKdoctoragentParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}