// Code generated by go-swagger; DO NOT EDIT.

// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package echo

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// NewPutKdoctoragentParams creates a new PutKdoctoragentParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewPutKdoctoragentParams() *PutKdoctoragentParams {
	return &PutKdoctoragentParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewPutKdoctoragentParamsWithTimeout creates a new PutKdoctoragentParams object
// with the ability to set a timeout on a request.
func NewPutKdoctoragentParamsWithTimeout(timeout time.Duration) *PutKdoctoragentParams {
	return &PutKdoctoragentParams{
		timeout: timeout,
	}
}

// NewPutKdoctoragentParamsWithContext creates a new PutKdoctoragentParams object
// with the ability to set a context for a request.
func NewPutKdoctoragentParamsWithContext(ctx context.Context) *PutKdoctoragentParams {
	return &PutKdoctoragentParams{
		Context: ctx,
	}
}

// NewPutKdoctoragentParamsWithHTTPClient creates a new PutKdoctoragentParams object
// with the ability to set a custom HTTPClient for a request.
func NewPutKdoctoragentParamsWithHTTPClient(client *http.Client) *PutKdoctoragentParams {
	return &PutKdoctoragentParams{
		HTTPClient: client,
	}
}

/*
PutKdoctoragentParams contains all the parameters to send to the API endpoint

	for the put kdoctoragent operation.

	Typically these are written to a http.Request.
*/
type PutKdoctoragentParams struct {

	/* Delay.

	   delay some second return response
	*/
	Delay *int64

	/* Task.

	   task name
	*/
	Task *string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the put kdoctoragent params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *PutKdoctoragentParams) WithDefaults() *PutKdoctoragentParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the put kdoctoragent params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *PutKdoctoragentParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the put kdoctoragent params
func (o *PutKdoctoragentParams) WithTimeout(timeout time.Duration) *PutKdoctoragentParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the put kdoctoragent params
func (o *PutKdoctoragentParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the put kdoctoragent params
func (o *PutKdoctoragentParams) WithContext(ctx context.Context) *PutKdoctoragentParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the put kdoctoragent params
func (o *PutKdoctoragentParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the put kdoctoragent params
func (o *PutKdoctoragentParams) WithHTTPClient(client *http.Client) *PutKdoctoragentParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the put kdoctoragent params
func (o *PutKdoctoragentParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithDelay adds the delay to the put kdoctoragent params
func (o *PutKdoctoragentParams) WithDelay(delay *int64) *PutKdoctoragentParams {
	o.SetDelay(delay)
	return o
}

// SetDelay adds the delay to the put kdoctoragent params
func (o *PutKdoctoragentParams) SetDelay(delay *int64) {
	o.Delay = delay
}

// WithTask adds the task to the put kdoctoragent params
func (o *PutKdoctoragentParams) WithTask(task *string) *PutKdoctoragentParams {
	o.SetTask(task)
	return o
}

// SetTask adds the task to the put kdoctoragent params
func (o *PutKdoctoragentParams) SetTask(task *string) {
	o.Task = task
}

// WriteToRequest writes these params to a swagger request
func (o *PutKdoctoragentParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.Delay != nil {

		// query param delay
		var qrDelay int64

		if o.Delay != nil {
			qrDelay = *o.Delay
		}
		qDelay := swag.FormatInt64(qrDelay)
		if qDelay != "" {

			if err := r.SetQueryParam("delay", qDelay); err != nil {
				return err
			}
		}
	}

	if o.Task != nil {

		// query param task
		var qrTask string

		if o.Task != nil {
			qrTask = *o.Task
		}
		qTask := qrTask
		if qTask != "" {

			if err := r.SetQueryParam("task", qTask); err != nil {
				return err
			}
		}
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
