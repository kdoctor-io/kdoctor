// Code generated by go-swagger; DO NOT EDIT.

// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package echo

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
)

// DeleteKdoctoragentReader is a Reader for the DeleteKdoctoragent structure.
type DeleteKdoctoragentReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *DeleteKdoctoragentReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewDeleteKdoctoragentOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 500:
		result := NewDeleteKdoctoragentInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		return nil, runtime.NewAPIError("response status code does not match any response statuses defined for this endpoint in the swagger spec", response, response.Code())
	}
}

// NewDeleteKdoctoragentOK creates a DeleteKdoctoragentOK with default headers values
func NewDeleteKdoctoragentOK() *DeleteKdoctoragentOK {
	return &DeleteKdoctoragentOK{}
}

/*
DeleteKdoctoragentOK describes a response with status code 200, with default header values.

Success
*/
type DeleteKdoctoragentOK struct {
}

// IsSuccess returns true when this delete kdoctoragent o k response has a 2xx status code
func (o *DeleteKdoctoragentOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this delete kdoctoragent o k response has a 3xx status code
func (o *DeleteKdoctoragentOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this delete kdoctoragent o k response has a 4xx status code
func (o *DeleteKdoctoragentOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this delete kdoctoragent o k response has a 5xx status code
func (o *DeleteKdoctoragentOK) IsServerError() bool {
	return false
}

// IsCode returns true when this delete kdoctoragent o k response a status code equal to that given
func (o *DeleteKdoctoragentOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the delete kdoctoragent o k response
func (o *DeleteKdoctoragentOK) Code() int {
	return 200
}

func (o *DeleteKdoctoragentOK) Error() string {
	return fmt.Sprintf("[DELETE /kdoctoragent][%d] deleteKdoctoragentOK ", 200)
}

func (o *DeleteKdoctoragentOK) String() string {
	return fmt.Sprintf("[DELETE /kdoctoragent][%d] deleteKdoctoragentOK ", 200)
}

func (o *DeleteKdoctoragentOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewDeleteKdoctoragentInternalServerError creates a DeleteKdoctoragentInternalServerError with default headers values
func NewDeleteKdoctoragentInternalServerError() *DeleteKdoctoragentInternalServerError {
	return &DeleteKdoctoragentInternalServerError{}
}

/*
DeleteKdoctoragentInternalServerError describes a response with status code 500, with default header values.

Failed
*/
type DeleteKdoctoragentInternalServerError struct {
}

// IsSuccess returns true when this delete kdoctoragent internal server error response has a 2xx status code
func (o *DeleteKdoctoragentInternalServerError) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this delete kdoctoragent internal server error response has a 3xx status code
func (o *DeleteKdoctoragentInternalServerError) IsRedirect() bool {
	return false
}

// IsClientError returns true when this delete kdoctoragent internal server error response has a 4xx status code
func (o *DeleteKdoctoragentInternalServerError) IsClientError() bool {
	return false
}

// IsServerError returns true when this delete kdoctoragent internal server error response has a 5xx status code
func (o *DeleteKdoctoragentInternalServerError) IsServerError() bool {
	return true
}

// IsCode returns true when this delete kdoctoragent internal server error response a status code equal to that given
func (o *DeleteKdoctoragentInternalServerError) IsCode(code int) bool {
	return code == 500
}

// Code gets the status code for the delete kdoctoragent internal server error response
func (o *DeleteKdoctoragentInternalServerError) Code() int {
	return 500
}

func (o *DeleteKdoctoragentInternalServerError) Error() string {
	return fmt.Sprintf("[DELETE /kdoctoragent][%d] deleteKdoctoragentInternalServerError ", 500)
}

func (o *DeleteKdoctoragentInternalServerError) String() string {
	return fmt.Sprintf("[DELETE /kdoctoragent][%d] deleteKdoctoragentInternalServerError ", 500)
}

func (o *DeleteKdoctoragentInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}
