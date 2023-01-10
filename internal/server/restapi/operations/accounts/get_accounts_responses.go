// Code generated by go-swagger; DO NOT EDIT.

package accounts

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"
)

// GetAccountsOKCode is the HTTP code returned for type GetAccountsOK
const GetAccountsOKCode int = 200

/*
GetAccountsOK OK

swagger:response getAccountsOK
*/
type GetAccountsOK struct {

	/*
	  In: Body
	*/
	Payload *GetAccountsOKBody `json:"body,omitempty"`
}

// NewGetAccountsOK creates GetAccountsOK with default headers values
func NewGetAccountsOK() *GetAccountsOK {

	return &GetAccountsOK{}
}

// WithPayload adds the payload to the get accounts o k response
func (o *GetAccountsOK) WithPayload(payload *GetAccountsOKBody) *GetAccountsOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get accounts o k response
func (o *GetAccountsOK) SetPayload(payload *GetAccountsOKBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetAccountsOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetAccountsInternalServerErrorCode is the HTTP code returned for type GetAccountsInternalServerError
const GetAccountsInternalServerErrorCode int = 500

/*
GetAccountsInternalServerError Ошибка на стороне сервера

swagger:response getAccountsInternalServerError
*/
type GetAccountsInternalServerError struct {

	/*
	  In: Body
	*/
	Payload interface{} `json:"body,omitempty"`
}

// NewGetAccountsInternalServerError creates GetAccountsInternalServerError with default headers values
func NewGetAccountsInternalServerError() *GetAccountsInternalServerError {

	return &GetAccountsInternalServerError{}
}

// WithPayload adds the payload to the get accounts internal server error response
func (o *GetAccountsInternalServerError) WithPayload(payload interface{}) *GetAccountsInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get accounts internal server error response
func (o *GetAccountsInternalServerError) SetPayload(payload interface{}) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetAccountsInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}