// Code generated by go-swagger; DO NOT EDIT.

package torrents

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

// SearchTorrentsReader is a Reader for the SearchTorrents structure.
type SearchTorrentsReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *SearchTorrentsReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewSearchTorrentsOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 500:
		result := NewSearchTorrentsInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		return nil, runtime.NewAPIError("response status code does not match any response statuses defined for this endpoint in the swagger spec", response, response.Code())
	}
}

// NewSearchTorrentsOK creates a SearchTorrentsOK with default headers values
func NewSearchTorrentsOK() *SearchTorrentsOK {
	return &SearchTorrentsOK{}
}

/*
SearchTorrentsOK describes a response with status code 200, with default header values.

OK
*/
type SearchTorrentsOK struct {
	Payload *SearchTorrentsOKBody
}

// IsSuccess returns true when this search torrents o k response has a 2xx status code
func (o *SearchTorrentsOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this search torrents o k response has a 3xx status code
func (o *SearchTorrentsOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this search torrents o k response has a 4xx status code
func (o *SearchTorrentsOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this search torrents o k response has a 5xx status code
func (o *SearchTorrentsOK) IsServerError() bool {
	return false
}

// IsCode returns true when this search torrents o k response a status code equal to that given
func (o *SearchTorrentsOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the search torrents o k response
func (o *SearchTorrentsOK) Code() int {
	return 200
}

func (o *SearchTorrentsOK) Error() string {
	return fmt.Sprintf("[GET /torrents/search][%d] searchTorrentsOK  %+v", 200, o.Payload)
}

func (o *SearchTorrentsOK) String() string {
	return fmt.Sprintf("[GET /torrents/search][%d] searchTorrentsOK  %+v", 200, o.Payload)
}

func (o *SearchTorrentsOK) GetPayload() *SearchTorrentsOKBody {
	return o.Payload
}

func (o *SearchTorrentsOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(SearchTorrentsOKBody)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewSearchTorrentsInternalServerError creates a SearchTorrentsInternalServerError with default headers values
func NewSearchTorrentsInternalServerError() *SearchTorrentsInternalServerError {
	return &SearchTorrentsInternalServerError{}
}

/*
SearchTorrentsInternalServerError describes a response with status code 500, with default header values.

Ошибка на стороне сервера
*/
type SearchTorrentsInternalServerError struct {
}

// IsSuccess returns true when this search torrents internal server error response has a 2xx status code
func (o *SearchTorrentsInternalServerError) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this search torrents internal server error response has a 3xx status code
func (o *SearchTorrentsInternalServerError) IsRedirect() bool {
	return false
}

// IsClientError returns true when this search torrents internal server error response has a 4xx status code
func (o *SearchTorrentsInternalServerError) IsClientError() bool {
	return false
}

// IsServerError returns true when this search torrents internal server error response has a 5xx status code
func (o *SearchTorrentsInternalServerError) IsServerError() bool {
	return true
}

// IsCode returns true when this search torrents internal server error response a status code equal to that given
func (o *SearchTorrentsInternalServerError) IsCode(code int) bool {
	return code == 500
}

// Code gets the status code for the search torrents internal server error response
func (o *SearchTorrentsInternalServerError) Code() int {
	return 500
}

func (o *SearchTorrentsInternalServerError) Error() string {
	return fmt.Sprintf("[GET /torrents/search][%d] searchTorrentsInternalServerError ", 500)
}

func (o *SearchTorrentsInternalServerError) String() string {
	return fmt.Sprintf("[GET /torrents/search][%d] searchTorrentsInternalServerError ", 500)
}

func (o *SearchTorrentsInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

/*
SearchTorrentsOKBody search torrents o k body
swagger:model SearchTorrentsOKBody
*/
type SearchTorrentsOKBody struct {

	// results
	// Required: true
	Results []*models.SearchTorrentsResult `json:"results"`
}

// Validate validates this search torrents o k body
func (o *SearchTorrentsOKBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateResults(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *SearchTorrentsOKBody) validateResults(formats strfmt.Registry) error {

	if err := validate.Required("searchTorrentsOK"+"."+"results", "body", o.Results); err != nil {
		return err
	}

	for i := 0; i < len(o.Results); i++ {
		if swag.IsZero(o.Results[i]) { // not required
			continue
		}

		if o.Results[i] != nil {
			if err := o.Results[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("searchTorrentsOK" + "." + "results" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("searchTorrentsOK" + "." + "results" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// ContextValidate validate this search torrents o k body based on the context it is used
func (o *SearchTorrentsOKBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := o.contextValidateResults(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *SearchTorrentsOKBody) contextValidateResults(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(o.Results); i++ {

		if o.Results[i] != nil {
			if err := o.Results[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("searchTorrentsOK" + "." + "results" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("searchTorrentsOK" + "." + "results" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (o *SearchTorrentsOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *SearchTorrentsOKBody) UnmarshalBinary(b []byte) error {
	var res SearchTorrentsOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
