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

// DeleteKdoctoragentOKCode is the HTTP code returned for type DeleteKdoctoragentOK
const DeleteKdoctoragentOKCode int = 200

/*
DeleteKdoctoragentOK Success

swagger:response deleteKdoctoragentOK
*/
type DeleteKdoctoragentOK struct {
}

// NewDeleteKdoctoragentOK creates DeleteKdoctoragentOK with default headers values
func NewDeleteKdoctoragentOK() *DeleteKdoctoragentOK {

	return &DeleteKdoctoragentOK{}
}

// WriteResponse to the client
func (o *DeleteKdoctoragentOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(200)
}

// DeleteKdoctoragentInternalServerErrorCode is the HTTP code returned for type DeleteKdoctoragentInternalServerError
const DeleteKdoctoragentInternalServerErrorCode int = 500

/*
DeleteKdoctoragentInternalServerError Failed

swagger:response deleteKdoctoragentInternalServerError
*/
type DeleteKdoctoragentInternalServerError struct {
}

// NewDeleteKdoctoragentInternalServerError creates DeleteKdoctoragentInternalServerError with default headers values
func NewDeleteKdoctoragentInternalServerError() *DeleteKdoctoragentInternalServerError {

	return &DeleteKdoctoragentInternalServerError{}
}

// WriteResponse to the client
func (o *DeleteKdoctoragentInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(500)
}