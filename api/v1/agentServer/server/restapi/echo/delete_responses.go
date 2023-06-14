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

// DeleteOKCode is the HTTP code returned for type DeleteOK
const DeleteOKCode int = 200

/*
DeleteOK Success

swagger:response deleteOK
*/
type DeleteOK struct {
}

// NewDeleteOK creates DeleteOK with default headers values
func NewDeleteOK() *DeleteOK {

	return &DeleteOK{}
}

// WriteResponse to the client
func (o *DeleteOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(200)
}

// DeleteInternalServerErrorCode is the HTTP code returned for type DeleteInternalServerError
const DeleteInternalServerErrorCode int = 500

/*
DeleteInternalServerError Failed

swagger:response deleteInternalServerError
*/
type DeleteInternalServerError struct {
}

// NewDeleteInternalServerError creates DeleteInternalServerError with default headers values
func NewDeleteInternalServerError() *DeleteInternalServerError {

	return &DeleteInternalServerError{}
}

// WriteResponse to the client
func (o *DeleteInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(500)
}
