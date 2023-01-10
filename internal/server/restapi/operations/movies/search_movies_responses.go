// Code generated by go-swagger; DO NOT EDIT.

package movies

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"
)

// SearchMoviesOKCode is the HTTP code returned for type SearchMoviesOK
const SearchMoviesOKCode int = 200

/*
SearchMoviesOK OK

swagger:response searchMoviesOK
*/
type SearchMoviesOK struct {

	/*
	  In: Body
	*/
	Payload *SearchMoviesOKBody `json:"body,omitempty"`
}

// NewSearchMoviesOK creates SearchMoviesOK with default headers values
func NewSearchMoviesOK() *SearchMoviesOK {

	return &SearchMoviesOK{}
}

// WithPayload adds the payload to the search movies o k response
func (o *SearchMoviesOK) WithPayload(payload *SearchMoviesOKBody) *SearchMoviesOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the search movies o k response
func (o *SearchMoviesOK) SetPayload(payload *SearchMoviesOKBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *SearchMoviesOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// SearchMoviesInternalServerErrorCode is the HTTP code returned for type SearchMoviesInternalServerError
const SearchMoviesInternalServerErrorCode int = 500

/*
SearchMoviesInternalServerError Ошибка на стороне сервера

swagger:response searchMoviesInternalServerError
*/
type SearchMoviesInternalServerError struct {

	/*
	  In: Body
	*/
	Payload interface{} `json:"body,omitempty"`
}

// NewSearchMoviesInternalServerError creates SearchMoviesInternalServerError with default headers values
func NewSearchMoviesInternalServerError() *SearchMoviesInternalServerError {

	return &SearchMoviesInternalServerError{}
}

// WithPayload adds the payload to the search movies internal server error response
func (o *SearchMoviesInternalServerError) WithPayload(payload interface{}) *SearchMoviesInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the search movies internal server error response
func (o *SearchMoviesInternalServerError) SetPayload(payload interface{}) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *SearchMoviesInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}