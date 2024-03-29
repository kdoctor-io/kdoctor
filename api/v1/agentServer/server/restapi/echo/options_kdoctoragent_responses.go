// Code generated by go-swagger; DO NOT EDIT.

// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package echo

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"
)

// OptionsKdoctoragentOKCode is the HTTP code returned for type OptionsKdoctoragentOK
const OptionsKdoctoragentOKCode int = 200

/*
OptionsKdoctoragentOK Success

swagger:response optionsKdoctoragentOK
*/
type OptionsKdoctoragentOK struct {
}

// NewOptionsKdoctoragentOK creates OptionsKdoctoragentOK with default headers values
func NewOptionsKdoctoragentOK() *OptionsKdoctoragentOK {

	return &OptionsKdoctoragentOK{}
}

// WriteResponse to the client
func (o *OptionsKdoctoragentOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(200)
}

// OptionsKdoctoragentInternalServerErrorCode is the HTTP code returned for type OptionsKdoctoragentInternalServerError
const OptionsKdoctoragentInternalServerErrorCode int = 500

/*
OptionsKdoctoragentInternalServerError Failed

swagger:response optionsKdoctoragentInternalServerError
*/
type OptionsKdoctoragentInternalServerError struct {
}

// NewOptionsKdoctoragentInternalServerError creates OptionsKdoctoragentInternalServerError with default headers values
func NewOptionsKdoctoragentInternalServerError() *OptionsKdoctoragentInternalServerError {

	return &OptionsKdoctoragentInternalServerError{}
}

// WriteResponse to the client
func (o *OptionsKdoctoragentInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(500)
}
