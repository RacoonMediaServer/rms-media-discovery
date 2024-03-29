// Code generated by go-swagger; DO NOT EDIT.

package torrents

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"
)

// SearchTorrentsAsyncOKCode is the HTTP code returned for type SearchTorrentsAsyncOK
const SearchTorrentsAsyncOKCode int = 200

/*
SearchTorrentsAsyncOK OK

swagger:response searchTorrentsAsyncOK
*/
type SearchTorrentsAsyncOK struct {

	/*
	  In: Body
	*/
	Payload *SearchTorrentsAsyncOKBody `json:"body,omitempty"`
}

// NewSearchTorrentsAsyncOK creates SearchTorrentsAsyncOK with default headers values
func NewSearchTorrentsAsyncOK() *SearchTorrentsAsyncOK {

	return &SearchTorrentsAsyncOK{}
}

// WithPayload adds the payload to the search torrents async o k response
func (o *SearchTorrentsAsyncOK) WithPayload(payload *SearchTorrentsAsyncOKBody) *SearchTorrentsAsyncOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the search torrents async o k response
func (o *SearchTorrentsAsyncOK) SetPayload(payload *SearchTorrentsAsyncOKBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *SearchTorrentsAsyncOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// SearchTorrentsAsyncInternalServerErrorCode is the HTTP code returned for type SearchTorrentsAsyncInternalServerError
const SearchTorrentsAsyncInternalServerErrorCode int = 500

/*
SearchTorrentsAsyncInternalServerError Ошибка на стороне сервера

swagger:response searchTorrentsAsyncInternalServerError
*/
type SearchTorrentsAsyncInternalServerError struct {
}

// NewSearchTorrentsAsyncInternalServerError creates SearchTorrentsAsyncInternalServerError with default headers values
func NewSearchTorrentsAsyncInternalServerError() *SearchTorrentsAsyncInternalServerError {

	return &SearchTorrentsAsyncInternalServerError{}
}

// WriteResponse to the client
func (o *SearchTorrentsAsyncInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(500)
}
