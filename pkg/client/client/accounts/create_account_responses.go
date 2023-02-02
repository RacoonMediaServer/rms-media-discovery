// Code generated by go-swagger; DO NOT EDIT.

package accounts

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"fmt"
	"io"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// CreateAccountReader is a Reader for the CreateAccount structure.
type CreateAccountReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *CreateAccountReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewCreateAccountOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 500:
		result := NewCreateAccountInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		return nil, runtime.NewAPIError("response status code does not match any response statuses defined for this endpoint in the swagger spec", response, response.Code())
	}
}

// NewCreateAccountOK creates a CreateAccountOK with default headers values
func NewCreateAccountOK() *CreateAccountOK {
	return &CreateAccountOK{}
}

/*
CreateAccountOK describes a response with status code 200, with default header values.

OK
*/
type CreateAccountOK struct {
	Payload *CreateAccountOKBody
}

// IsSuccess returns true when this create account o k response has a 2xx status code
func (o *CreateAccountOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this create account o k response has a 3xx status code
func (o *CreateAccountOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this create account o k response has a 4xx status code
func (o *CreateAccountOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this create account o k response has a 5xx status code
func (o *CreateAccountOK) IsServerError() bool {
	return false
}

// IsCode returns true when this create account o k response a status code equal to that given
func (o *CreateAccountOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the create account o k response
func (o *CreateAccountOK) Code() int {
	return 200
}

func (o *CreateAccountOK) Error() string {
	return fmt.Sprintf("[POST /admin/accounts][%d] createAccountOK  %+v", 200, o.Payload)
}

func (o *CreateAccountOK) String() string {
	return fmt.Sprintf("[POST /admin/accounts][%d] createAccountOK  %+v", 200, o.Payload)
}

func (o *CreateAccountOK) GetPayload() *CreateAccountOKBody {
	return o.Payload
}

func (o *CreateAccountOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(CreateAccountOKBody)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewCreateAccountInternalServerError creates a CreateAccountInternalServerError with default headers values
func NewCreateAccountInternalServerError() *CreateAccountInternalServerError {
	return &CreateAccountInternalServerError{}
}

/*
CreateAccountInternalServerError describes a response with status code 500, with default header values.

Ошибка на стороне сервера
*/
type CreateAccountInternalServerError struct {
}

// IsSuccess returns true when this create account internal server error response has a 2xx status code
func (o *CreateAccountInternalServerError) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this create account internal server error response has a 3xx status code
func (o *CreateAccountInternalServerError) IsRedirect() bool {
	return false
}

// IsClientError returns true when this create account internal server error response has a 4xx status code
func (o *CreateAccountInternalServerError) IsClientError() bool {
	return false
}

// IsServerError returns true when this create account internal server error response has a 5xx status code
func (o *CreateAccountInternalServerError) IsServerError() bool {
	return true
}

// IsCode returns true when this create account internal server error response a status code equal to that given
func (o *CreateAccountInternalServerError) IsCode(code int) bool {
	return code == 500
}

// Code gets the status code for the create account internal server error response
func (o *CreateAccountInternalServerError) Code() int {
	return 500
}

func (o *CreateAccountInternalServerError) Error() string {
	return fmt.Sprintf("[POST /admin/accounts][%d] createAccountInternalServerError ", 500)
}

func (o *CreateAccountInternalServerError) String() string {
	return fmt.Sprintf("[POST /admin/accounts][%d] createAccountInternalServerError ", 500)
}

func (o *CreateAccountInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

/*
CreateAccountOKBody create account o k body
swagger:model CreateAccountOKBody
*/
type CreateAccountOKBody struct {

	// id
	// Required: true
	ID *string `json:"id"`
}

// Validate validates this create account o k body
func (o *CreateAccountOKBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateID(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *CreateAccountOKBody) validateID(formats strfmt.Registry) error {

	if err := validate.Required("createAccountOK"+"."+"id", "body", o.ID); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this create account o k body based on context it is used
func (o *CreateAccountOKBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *CreateAccountOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *CreateAccountOKBody) UnmarshalBinary(b []byte) error {
	var res CreateAccountOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
