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

// PostKdoctoragentOKCode is the HTTP code returned for type PostKdoctoragentOK
const PostKdoctoragentOKCode int = 200

/*
PostKdoctoragentOK Success

swagger:response postKdoctoragentOK
*/
type PostKdoctoragentOK struct {
}

// NewPostKdoctoragentOK creates PostKdoctoragentOK with default headers values
func NewPostKdoctoragentOK() *PostKdoctoragentOK {

	return &PostKdoctoragentOK{}
}

// WriteResponse to the client
func (o *PostKdoctoragentOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(200)
}

// PostKdoctoragentInternalServerErrorCode is the HTTP code returned for type PostKdoctoragentInternalServerError
const PostKdoctoragentInternalServerErrorCode int = 500

/*
PostKdoctoragentInternalServerError Failed

swagger:response postKdoctoragentInternalServerError
*/
type PostKdoctoragentInternalServerError struct {
}

// NewPostKdoctoragentInternalServerError creates PostKdoctoragentInternalServerError with default headers values
func NewPostKdoctoragentInternalServerError() *PostKdoctoragentInternalServerError {

	return &PostKdoctoragentInternalServerError{}
}

// WriteResponse to the client
func (o *PostKdoctoragentInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(500)
}
