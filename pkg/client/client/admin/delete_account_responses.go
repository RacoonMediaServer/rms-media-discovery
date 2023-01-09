// Code generated by go-swagger; DO NOT EDIT.

package admin

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
)

// DeleteAccountReader is a Reader for the DeleteAccount structure.
type DeleteAccountReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *DeleteAccountReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewDeleteAccountOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 404:
		result := NewDeleteAccountNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewDeleteAccountInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		return nil, runtime.NewAPIError("response status code does not match any response statuses defined for this endpoint in the swagger spec", response, response.Code())
	}
}

// NewDeleteAccountOK creates a DeleteAccountOK with default headers values
func NewDeleteAccountOK() *DeleteAccountOK {
	return &DeleteAccountOK{}
}

/*
DeleteAccountOK describes a response with status code 200, with default header values.

OK
*/
type DeleteAccountOK struct {
}

// IsSuccess returns true when this delete account o k response has a 2xx status code
func (o *DeleteAccountOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this delete account o k response has a 3xx status code
func (o *DeleteAccountOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this delete account o k response has a 4xx status code
func (o *DeleteAccountOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this delete account o k response has a 5xx status code
func (o *DeleteAccountOK) IsServerError() bool {
	return false
}

// IsCode returns true when this delete account o k response a status code equal to that given
func (o *DeleteAccountOK) IsCode(code int) bool {
	return code == 200
}

func (o *DeleteAccountOK) Error() string {
	return fmt.Sprintf("[DELETE /admin/accounts/{id}][%d] deleteAccountOK ", 200)
}

func (o *DeleteAccountOK) String() string {
	return fmt.Sprintf("[DELETE /admin/accounts/{id}][%d] deleteAccountOK ", 200)
}

func (o *DeleteAccountOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewDeleteAccountNotFound creates a DeleteAccountNotFound with default headers values
func NewDeleteAccountNotFound() *DeleteAccountNotFound {
	return &DeleteAccountNotFound{}
}

/*
DeleteAccountNotFound describes a response with status code 404, with default header values.

Аккаунт не найден
*/
type DeleteAccountNotFound struct {
}

// IsSuccess returns true when this delete account not found response has a 2xx status code
func (o *DeleteAccountNotFound) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this delete account not found response has a 3xx status code
func (o *DeleteAccountNotFound) IsRedirect() bool {
	return false
}

// IsClientError returns true when this delete account not found response has a 4xx status code
func (o *DeleteAccountNotFound) IsClientError() bool {
	return true
}

// IsServerError returns true when this delete account not found response has a 5xx status code
func (o *DeleteAccountNotFound) IsServerError() bool {
	return false
}

// IsCode returns true when this delete account not found response a status code equal to that given
func (o *DeleteAccountNotFound) IsCode(code int) bool {
	return code == 404
}

func (o *DeleteAccountNotFound) Error() string {
	return fmt.Sprintf("[DELETE /admin/accounts/{id}][%d] deleteAccountNotFound ", 404)
}

func (o *DeleteAccountNotFound) String() string {
	return fmt.Sprintf("[DELETE /admin/accounts/{id}][%d] deleteAccountNotFound ", 404)
}

func (o *DeleteAccountNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewDeleteAccountInternalServerError creates a DeleteAccountInternalServerError with default headers values
func NewDeleteAccountInternalServerError() *DeleteAccountInternalServerError {
	return &DeleteAccountInternalServerError{}
}

/*
DeleteAccountInternalServerError describes a response with status code 500, with default header values.

Ошибка на стороне сервера
*/
type DeleteAccountInternalServerError struct {
	Payload interface{}
}

// IsSuccess returns true when this delete account internal server error response has a 2xx status code
func (o *DeleteAccountInternalServerError) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this delete account internal server error response has a 3xx status code
func (o *DeleteAccountInternalServerError) IsRedirect() bool {
	return false
}

// IsClientError returns true when this delete account internal server error response has a 4xx status code
func (o *DeleteAccountInternalServerError) IsClientError() bool {
	return false
}

// IsServerError returns true when this delete account internal server error response has a 5xx status code
func (o *DeleteAccountInternalServerError) IsServerError() bool {
	return true
}

// IsCode returns true when this delete account internal server error response a status code equal to that given
func (o *DeleteAccountInternalServerError) IsCode(code int) bool {
	return code == 500
}

func (o *DeleteAccountInternalServerError) Error() string {
	return fmt.Sprintf("[DELETE /admin/accounts/{id}][%d] deleteAccountInternalServerError  %+v", 500, o.Payload)
}

func (o *DeleteAccountInternalServerError) String() string {
	return fmt.Sprintf("[DELETE /admin/accounts/{id}][%d] deleteAccountInternalServerError  %+v", 500, o.Payload)
}

func (o *DeleteAccountInternalServerError) GetPayload() interface{} {
	return o.Payload
}

func (o *DeleteAccountInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
