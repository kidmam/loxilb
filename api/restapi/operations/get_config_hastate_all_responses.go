// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/loxilb-io/loxilb/api/models"
)

// GetConfigHastateAllOKCode is the HTTP code returned for type GetConfigHastateAllOK
const GetConfigHastateAllOKCode int = 200

/*
GetConfigHastateAllOK OK

swagger:response getConfigHastateAllOK
*/
type GetConfigHastateAllOK struct {

	/*
	  In: Body
	*/
	Payload *GetConfigHastateAllOKBody `json:"body,omitempty"`
}

// NewGetConfigHastateAllOK creates GetConfigHastateAllOK with default headers values
func NewGetConfigHastateAllOK() *GetConfigHastateAllOK {

	return &GetConfigHastateAllOK{}
}

// WithPayload adds the payload to the get config hastate all o k response
func (o *GetConfigHastateAllOK) WithPayload(payload *GetConfigHastateAllOKBody) *GetConfigHastateAllOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get config hastate all o k response
func (o *GetConfigHastateAllOK) SetPayload(payload *GetConfigHastateAllOKBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetConfigHastateAllOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetConfigHastateAllUnauthorizedCode is the HTTP code returned for type GetConfigHastateAllUnauthorized
const GetConfigHastateAllUnauthorizedCode int = 401

/*
GetConfigHastateAllUnauthorized Invalid authentication credentials

swagger:response getConfigHastateAllUnauthorized
*/
type GetConfigHastateAllUnauthorized struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetConfigHastateAllUnauthorized creates GetConfigHastateAllUnauthorized with default headers values
func NewGetConfigHastateAllUnauthorized() *GetConfigHastateAllUnauthorized {

	return &GetConfigHastateAllUnauthorized{}
}

// WithPayload adds the payload to the get config hastate all unauthorized response
func (o *GetConfigHastateAllUnauthorized) WithPayload(payload *models.Error) *GetConfigHastateAllUnauthorized {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get config hastate all unauthorized response
func (o *GetConfigHastateAllUnauthorized) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetConfigHastateAllUnauthorized) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(401)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetConfigHastateAllInternalServerErrorCode is the HTTP code returned for type GetConfigHastateAllInternalServerError
const GetConfigHastateAllInternalServerErrorCode int = 500

/*
GetConfigHastateAllInternalServerError Internal service error

swagger:response getConfigHastateAllInternalServerError
*/
type GetConfigHastateAllInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetConfigHastateAllInternalServerError creates GetConfigHastateAllInternalServerError with default headers values
func NewGetConfigHastateAllInternalServerError() *GetConfigHastateAllInternalServerError {

	return &GetConfigHastateAllInternalServerError{}
}

// WithPayload adds the payload to the get config hastate all internal server error response
func (o *GetConfigHastateAllInternalServerError) WithPayload(payload *models.Error) *GetConfigHastateAllInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get config hastate all internal server error response
func (o *GetConfigHastateAllInternalServerError) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetConfigHastateAllInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetConfigHastateAllServiceUnavailableCode is the HTTP code returned for type GetConfigHastateAllServiceUnavailable
const GetConfigHastateAllServiceUnavailableCode int = 503

/*
GetConfigHastateAllServiceUnavailable Maintanence mode

swagger:response getConfigHastateAllServiceUnavailable
*/
type GetConfigHastateAllServiceUnavailable struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetConfigHastateAllServiceUnavailable creates GetConfigHastateAllServiceUnavailable with default headers values
func NewGetConfigHastateAllServiceUnavailable() *GetConfigHastateAllServiceUnavailable {

	return &GetConfigHastateAllServiceUnavailable{}
}

// WithPayload adds the payload to the get config hastate all service unavailable response
func (o *GetConfigHastateAllServiceUnavailable) WithPayload(payload *models.Error) *GetConfigHastateAllServiceUnavailable {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get config hastate all service unavailable response
func (o *GetConfigHastateAllServiceUnavailable) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetConfigHastateAllServiceUnavailable) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(503)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
