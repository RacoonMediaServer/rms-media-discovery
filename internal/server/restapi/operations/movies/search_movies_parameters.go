// Code generated by go-swagger; DO NOT EDIT.

package movies

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// NewSearchMoviesParams creates a new SearchMoviesParams object
//
// There are no default values defined in the spec.
func NewSearchMoviesParams() SearchMoviesParams {

	return SearchMoviesParams{}
}

// SearchMoviesParams contains all the bound params for the search movies operation
// typically these are obtained from a http.Request
//
// swagger:parameters searchMovies
type SearchMoviesParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*Ограничение на кол-во результатов
	  Minimum: 1
	  In: query
	*/
	Limit *int64
	/*Искомый запрос
	  Required: true
	  Max Length: 128
	  Min Length: 2
	  In: query
	*/
	Q string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewSearchMoviesParams() beforehand.
func (o *SearchMoviesParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	qs := runtime.Values(r.URL.Query())

	qLimit, qhkLimit, _ := qs.GetOK("limit")
	if err := o.bindLimit(qLimit, qhkLimit, route.Formats); err != nil {
		res = append(res, err)
	}

	qQ, qhkQ, _ := qs.GetOK("q")
	if err := o.bindQ(qQ, qhkQ, route.Formats); err != nil {
		res = append(res, err)
	}
	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindLimit binds and validates parameter Limit from query.
func (o *SearchMoviesParams) bindLimit(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false

	if raw == "" { // empty values pass all other validations
		return nil
	}

	value, err := swag.ConvertInt64(raw)
	if err != nil {
		return errors.InvalidType("limit", "query", "int64", raw)
	}
	o.Limit = &value

	if err := o.validateLimit(formats); err != nil {
		return err
	}

	return nil
}

// validateLimit carries on validations for parameter Limit
func (o *SearchMoviesParams) validateLimit(formats strfmt.Registry) error {

	if err := validate.MinimumInt("limit", "query", *o.Limit, 1, false); err != nil {
		return err
	}

	return nil
}

// bindQ binds and validates parameter Q from query.
func (o *SearchMoviesParams) bindQ(rawData []string, hasKey bool, formats strfmt.Registry) error {
	if !hasKey {
		return errors.Required("q", "query", rawData)
	}
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// AllowEmptyValue: false

	if err := validate.RequiredString("q", "query", raw); err != nil {
		return err
	}
	o.Q = raw

	if err := o.validateQ(formats); err != nil {
		return err
	}

	return nil
}

// validateQ carries on validations for parameter Q
func (o *SearchMoviesParams) validateQ(formats strfmt.Registry) error {

	if err := validate.MinLength("q", "query", o.Q, 2); err != nil {
		return err
	}

	if err := validate.MaxLength("q", "query", o.Q, 128); err != nil {
		return err
	}

	return nil
}