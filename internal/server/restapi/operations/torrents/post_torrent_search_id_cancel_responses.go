// Code generated by go-swagger; DO NOT EDIT.

package torrents

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"
)

// PostTorrentSearchIDCancelOKCode is the HTTP code returned for type PostTorrentSearchIDCancelOK
const PostTorrentSearchIDCancelOKCode int = 200

/*
PostTorrentSearchIDCancelOK OK

swagger:response postTorrentSearchIdCancelOK
*/
type PostTorrentSearchIDCancelOK struct {
}

// NewPostTorrentSearchIDCancelOK creates PostTorrentSearchIDCancelOK with default headers values
func NewPostTorrentSearchIDCancelOK() *PostTorrentSearchIDCancelOK {

	return &PostTorrentSearchIDCancelOK{}
}

// WriteResponse to the client
func (o *PostTorrentSearchIDCancelOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(200)
}

// PostTorrentSearchIDCancelNotFoundCode is the HTTP code returned for type PostTorrentSearchIDCancelNotFound
const PostTorrentSearchIDCancelNotFoundCode int = 404

/*
PostTorrentSearchIDCancelNotFound Задача поиска не найдена

swagger:response postTorrentSearchIdCancelNotFound
*/
type PostTorrentSearchIDCancelNotFound struct {
}

// NewPostTorrentSearchIDCancelNotFound creates PostTorrentSearchIDCancelNotFound with default headers values
func NewPostTorrentSearchIDCancelNotFound() *PostTorrentSearchIDCancelNotFound {

	return &PostTorrentSearchIDCancelNotFound{}
}

// WriteResponse to the client
func (o *PostTorrentSearchIDCancelNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(404)
}

// PostTorrentSearchIDCancelInternalServerErrorCode is the HTTP code returned for type PostTorrentSearchIDCancelInternalServerError
const PostTorrentSearchIDCancelInternalServerErrorCode int = 500

/*
PostTorrentSearchIDCancelInternalServerError Ошибка на стороне сервера

swagger:response postTorrentSearchIdCancelInternalServerError
*/
type PostTorrentSearchIDCancelInternalServerError struct {
}

// NewPostTorrentSearchIDCancelInternalServerError creates PostTorrentSearchIDCancelInternalServerError with default headers values
func NewPostTorrentSearchIDCancelInternalServerError() *PostTorrentSearchIDCancelInternalServerError {

	return &PostTorrentSearchIDCancelInternalServerError{}
}

// WriteResponse to the client
func (o *PostTorrentSearchIDCancelInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(500)
}
