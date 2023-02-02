// Code generated by go-swagger; DO NOT EDIT.

package movies

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

	"git.rms.local/RacoonMediaServer/rms-media-discovery/pkg/client/models"
)

// SearchMoviesReader is a Reader for the SearchMovies structure.
type SearchMoviesReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *SearchMoviesReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewSearchMoviesOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 500:
		result := NewSearchMoviesInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		return nil, runtime.NewAPIError("response status code does not match any response statuses defined for this endpoint in the swagger spec", response, response.Code())
	}
}

// NewSearchMoviesOK creates a SearchMoviesOK with default headers values
func NewSearchMoviesOK() *SearchMoviesOK {
	return &SearchMoviesOK{}
}

/*
SearchMoviesOK describes a response with status code 200, with default header values.

OK
*/
type SearchMoviesOK struct {
	Payload *SearchMoviesOKBody
}

// IsSuccess returns true when this search movies o k response has a 2xx status code
func (o *SearchMoviesOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this search movies o k response has a 3xx status code
func (o *SearchMoviesOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this search movies o k response has a 4xx status code
func (o *SearchMoviesOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this search movies o k response has a 5xx status code
func (o *SearchMoviesOK) IsServerError() bool {
	return false
}

// IsCode returns true when this search movies o k response a status code equal to that given
func (o *SearchMoviesOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the search movies o k response
func (o *SearchMoviesOK) Code() int {
	return 200
}

func (o *SearchMoviesOK) Error() string {
	return fmt.Sprintf("[GET /movies/search][%d] searchMoviesOK  %+v", 200, o.Payload)
}

func (o *SearchMoviesOK) String() string {
	return fmt.Sprintf("[GET /movies/search][%d] searchMoviesOK  %+v", 200, o.Payload)
}

func (o *SearchMoviesOK) GetPayload() *SearchMoviesOKBody {
	return o.Payload
}

func (o *SearchMoviesOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(SearchMoviesOKBody)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewSearchMoviesInternalServerError creates a SearchMoviesInternalServerError with default headers values
func NewSearchMoviesInternalServerError() *SearchMoviesInternalServerError {
	return &SearchMoviesInternalServerError{}
}

/*
SearchMoviesInternalServerError describes a response with status code 500, with default header values.

Ошибка на стороне сервера
*/
type SearchMoviesInternalServerError struct {
}

// IsSuccess returns true when this search movies internal server error response has a 2xx status code
func (o *SearchMoviesInternalServerError) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this search movies internal server error response has a 3xx status code
func (o *SearchMoviesInternalServerError) IsRedirect() bool {
	return false
}

// IsClientError returns true when this search movies internal server error response has a 4xx status code
func (o *SearchMoviesInternalServerError) IsClientError() bool {
	return false
}

// IsServerError returns true when this search movies internal server error response has a 5xx status code
func (o *SearchMoviesInternalServerError) IsServerError() bool {
	return true
}

// IsCode returns true when this search movies internal server error response a status code equal to that given
func (o *SearchMoviesInternalServerError) IsCode(code int) bool {
	return code == 500
}

// Code gets the status code for the search movies internal server error response
func (o *SearchMoviesInternalServerError) Code() int {
	return 500
}

func (o *SearchMoviesInternalServerError) Error() string {
	return fmt.Sprintf("[GET /movies/search][%d] searchMoviesInternalServerError ", 500)
}

func (o *SearchMoviesInternalServerError) String() string {
	return fmt.Sprintf("[GET /movies/search][%d] searchMoviesInternalServerError ", 500)
}

func (o *SearchMoviesInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

/*
SearchMoviesOKBody search movies o k body
swagger:model SearchMoviesOKBody
*/
type SearchMoviesOKBody struct {

	// results
	// Required: true
	Results []*models.SearchMoviesResult `json:"results"`
}

// Validate validates this search movies o k body
func (o *SearchMoviesOKBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateResults(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *SearchMoviesOKBody) validateResults(formats strfmt.Registry) error {

	if err := validate.Required("searchMoviesOK"+"."+"results", "body", o.Results); err != nil {
		return err
	}

	for i := 0; i < len(o.Results); i++ {
		if swag.IsZero(o.Results[i]) { // not required
			continue
		}

		if o.Results[i] != nil {
			if err := o.Results[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("searchMoviesOK" + "." + "results" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("searchMoviesOK" + "." + "results" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// ContextValidate validate this search movies o k body based on the context it is used
func (o *SearchMoviesOKBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := o.contextValidateResults(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *SearchMoviesOKBody) contextValidateResults(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(o.Results); i++ {

		if o.Results[i] != nil {
			if err := o.Results[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("searchMoviesOK" + "." + "results" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("searchMoviesOK" + "." + "results" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (o *SearchMoviesOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *SearchMoviesOKBody) UnmarshalBinary(b []byte) error {
	var res SearchMoviesOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
