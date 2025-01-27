// Code generated by go-swagger; DO NOT EDIT.

package torrents

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
)

// DownloadTorrentReader is a Reader for the DownloadTorrent structure.
type DownloadTorrentReader struct {
	formats strfmt.Registry
	writer  io.Writer
}

// ReadResponse reads a server response into the received o.
func (o *DownloadTorrentReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewDownloadTorrentOK(o.writer)
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 404:
		result := NewDownloadTorrentNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewDownloadTorrentInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		return nil, runtime.NewAPIError("[GET /torrents/download] downloadTorrent", response, response.Code())
	}
}

// NewDownloadTorrentOK creates a DownloadTorrentOK with default headers values
func NewDownloadTorrentOK(writer io.Writer) *DownloadTorrentOK {
	return &DownloadTorrentOK{

		Payload: writer,
	}
}

/*
DownloadTorrentOK describes a response with status code 200, with default header values.

OK
*/
type DownloadTorrentOK struct {
	Payload io.Writer
}

// IsSuccess returns true when this download torrent o k response has a 2xx status code
func (o *DownloadTorrentOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this download torrent o k response has a 3xx status code
func (o *DownloadTorrentOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this download torrent o k response has a 4xx status code
func (o *DownloadTorrentOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this download torrent o k response has a 5xx status code
func (o *DownloadTorrentOK) IsServerError() bool {
	return false
}

// IsCode returns true when this download torrent o k response a status code equal to that given
func (o *DownloadTorrentOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the download torrent o k response
func (o *DownloadTorrentOK) Code() int {
	return 200
}

func (o *DownloadTorrentOK) Error() string {
	return fmt.Sprintf("[GET /torrents/download][%d] downloadTorrentOK", 200)
}

func (o *DownloadTorrentOK) String() string {
	return fmt.Sprintf("[GET /torrents/download][%d] downloadTorrentOK", 200)
}

func (o *DownloadTorrentOK) GetPayload() io.Writer {
	return o.Payload
}

func (o *DownloadTorrentOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewDownloadTorrentNotFound creates a DownloadTorrentNotFound with default headers values
func NewDownloadTorrentNotFound() *DownloadTorrentNotFound {
	return &DownloadTorrentNotFound{}
}

/*
DownloadTorrentNotFound describes a response with status code 404, with default header values.

Неверный хеш ссылки
*/
type DownloadTorrentNotFound struct {
}

// IsSuccess returns true when this download torrent not found response has a 2xx status code
func (o *DownloadTorrentNotFound) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this download torrent not found response has a 3xx status code
func (o *DownloadTorrentNotFound) IsRedirect() bool {
	return false
}

// IsClientError returns true when this download torrent not found response has a 4xx status code
func (o *DownloadTorrentNotFound) IsClientError() bool {
	return true
}

// IsServerError returns true when this download torrent not found response has a 5xx status code
func (o *DownloadTorrentNotFound) IsServerError() bool {
	return false
}

// IsCode returns true when this download torrent not found response a status code equal to that given
func (o *DownloadTorrentNotFound) IsCode(code int) bool {
	return code == 404
}

// Code gets the status code for the download torrent not found response
func (o *DownloadTorrentNotFound) Code() int {
	return 404
}

func (o *DownloadTorrentNotFound) Error() string {
	return fmt.Sprintf("[GET /torrents/download][%d] downloadTorrentNotFound", 404)
}

func (o *DownloadTorrentNotFound) String() string {
	return fmt.Sprintf("[GET /torrents/download][%d] downloadTorrentNotFound", 404)
}

func (o *DownloadTorrentNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewDownloadTorrentInternalServerError creates a DownloadTorrentInternalServerError with default headers values
func NewDownloadTorrentInternalServerError() *DownloadTorrentInternalServerError {
	return &DownloadTorrentInternalServerError{}
}

/*
DownloadTorrentInternalServerError describes a response with status code 500, with default header values.

Ошибка на стороне сервера
*/
type DownloadTorrentInternalServerError struct {
}

// IsSuccess returns true when this download torrent internal server error response has a 2xx status code
func (o *DownloadTorrentInternalServerError) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this download torrent internal server error response has a 3xx status code
func (o *DownloadTorrentInternalServerError) IsRedirect() bool {
	return false
}

// IsClientError returns true when this download torrent internal server error response has a 4xx status code
func (o *DownloadTorrentInternalServerError) IsClientError() bool {
	return false
}

// IsServerError returns true when this download torrent internal server error response has a 5xx status code
func (o *DownloadTorrentInternalServerError) IsServerError() bool {
	return true
}

// IsCode returns true when this download torrent internal server error response a status code equal to that given
func (o *DownloadTorrentInternalServerError) IsCode(code int) bool {
	return code == 500
}

// Code gets the status code for the download torrent internal server error response
func (o *DownloadTorrentInternalServerError) Code() int {
	return 500
}

func (o *DownloadTorrentInternalServerError) Error() string {
	return fmt.Sprintf("[GET /torrents/download][%d] downloadTorrentInternalServerError", 500)
}

func (o *DownloadTorrentInternalServerError) String() string {
	return fmt.Sprintf("[GET /torrents/download][%d] downloadTorrentInternalServerError", 500)
}

func (o *DownloadTorrentInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}
