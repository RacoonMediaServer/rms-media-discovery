// Code generated by go-swagger; DO NOT EDIT.

package accounts

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"fmt"
	"io"
	"strconv"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"

	"github.com/RacoonMediaServer/rms-media-discovery/pkg/client/models"
)

// GetAccountsReader is a Reader for the GetAccounts structure.
type GetAccountsReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetAccountsReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetAccountsOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 500:
		result := NewGetAccountsInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		return nil, runtime.NewAPIError("response status code does not match any response statuses defined for this endpoint in the swagger spec", response, response.Code())
	}
}

// NewGetAccountsOK creates a GetAccountsOK with default headers values
func NewGetAccountsOK() *GetAccountsOK {
	return &GetAccountsOK{}
}

/*
GetAccountsOK describes a response with status code 200, with default header values.

OK
*/
type GetAccountsOK struct {
	Payload *GetAccountsOKBody
}

// IsSuccess returns true when this get accounts o k response has a 2xx status code
func (o *GetAccountsOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this get accounts o k response has a 3xx status code
func (o *GetAccountsOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get accounts o k response has a 4xx status code
func (o *GetAccountsOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this get accounts o k response has a 5xx status code
func (o *GetAccountsOK) IsServerError() bool {
	return false
}

// IsCode returns true when this get accounts o k response a status code equal to that given
func (o *GetAccountsOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the get accounts o k response
func (o *GetAccountsOK) Code() int {
	return 200
}

func (o *GetAccountsOK) Error() string {
	return fmt.Sprintf("[GET /admin/accounts][%d] getAccountsOK  %+v", 200, o.Payload)
}

func (o *GetAccountsOK) String() string {
	return fmt.Sprintf("[GET /admin/accounts][%d] getAccountsOK  %+v", 200, o.Payload)
}

func (o *GetAccountsOK) GetPayload() *GetAccountsOKBody {
	return o.Payload
}

func (o *GetAccountsOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(GetAccountsOKBody)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetAccountsInternalServerError creates a GetAccountsInternalServerError with default headers values
func NewGetAccountsInternalServerError() *GetAccountsInternalServerError {
	return &GetAccountsInternalServerError{}
}

/*
GetAccountsInternalServerError describes a response with status code 500, with default header values.

Ошибка на стороне сервера
*/
type GetAccountsInternalServerError struct {
}

// IsSuccess returns true when this get accounts internal server error response has a 2xx status code
func (o *GetAccountsInternalServerError) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this get accounts internal server error response has a 3xx status code
func (o *GetAccountsInternalServerError) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get accounts internal server error response has a 4xx status code
func (o *GetAccountsInternalServerError) IsClientError() bool {
	return false
}

// IsServerError returns true when this get accounts internal server error response has a 5xx status code
func (o *GetAccountsInternalServerError) IsServerError() bool {
	return true
}

// IsCode returns true when this get accounts internal server error response a status code equal to that given
func (o *GetAccountsInternalServerError) IsCode(code int) bool {
	return code == 500
}

// Code gets the status code for the get accounts internal server error response
func (o *GetAccountsInternalServerError) Code() int {
	return 500
}

func (o *GetAccountsInternalServerError) Error() string {
	return fmt.Sprintf("[GET /admin/accounts][%d] getAccountsInternalServerError ", 500)
}

func (o *GetAccountsInternalServerError) String() string {
	return fmt.Sprintf("[GET /admin/accounts][%d] getAccountsInternalServerError ", 500)
}

func (o *GetAccountsInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

/*
GetAccountsOKBody get accounts o k body
swagger:model GetAccountsOKBody
*/
type GetAccountsOKBody struct {

	// results
	// Required: true
	Results []*models.Account `json:"results"`
}

// Validate validates this get accounts o k body
func (o *GetAccountsOKBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateResults(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *GetAccountsOKBody) validateResults(formats strfmt.Registry) error {

	if err := validate.Required("getAccountsOK"+"."+"results", "body", o.Results); err != nil {
		return err
	}

	for i := 0; i < len(o.Results); i++ {
		if swag.IsZero(o.Results[i]) { // not required
			continue
		}

		if o.Results[i] != nil {
			if err := o.Results[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("getAccountsOK" + "." + "results" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("getAccountsOK" + "." + "results" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// ContextValidate validate this get accounts o k body based on the context it is used
func (o *GetAccountsOKBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := o.contextValidateResults(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *GetAccountsOKBody) contextValidateResults(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(o.Results); i++ {

		if o.Results[i] != nil {
			if err := o.Results[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("getAccountsOK" + "." + "results" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("getAccountsOK" + "." + "results" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (o *GetAccountsOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *GetAccountsOKBody) UnmarshalBinary(b []byte) error {
	var res GetAccountsOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
