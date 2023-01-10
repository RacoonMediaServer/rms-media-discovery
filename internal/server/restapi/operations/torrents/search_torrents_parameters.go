// Code generated by go-swagger; DO NOT EDIT.

package torrents

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

// NewSearchTorrentsParams creates a new SearchTorrentsParams object
//
// There are no default values defined in the spec.
func NewSearchTorrentsParams() SearchTorrentsParams {

	return SearchTorrentsParams{}
}

// SearchTorrentsParams contains all the bound params for the search torrents operation
// typically these are obtained from a http.Request
//
// swagger:parameters searchTorrents
type SearchTorrentsParams struct {

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
	/*Подсказка, какого типа инфу искать
	  In: query
	*/
	Type *string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewSearchTorrentsParams() beforehand.
func (o *SearchTorrentsParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
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

	qType, qhkType, _ := qs.GetOK("type")
	if err := o.bindType(qType, qhkType, route.Formats); err != nil {
		res = append(res, err)
	}
	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindLimit binds and validates parameter Limit from query.
func (o *SearchTorrentsParams) bindLimit(rawData []string, hasKey bool, formats strfmt.Registry) error {
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
func (o *SearchTorrentsParams) validateLimit(formats strfmt.Registry) error {

	if err := validate.MinimumInt("limit", "query", *o.Limit, 1, false); err != nil {
		return err
	}

	return nil
}

// bindQ binds and validates parameter Q from query.
func (o *SearchTorrentsParams) bindQ(rawData []string, hasKey bool, formats strfmt.Registry) error {
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
func (o *SearchTorrentsParams) validateQ(formats strfmt.Registry) error {

	if err := validate.MinLength("q", "query", o.Q, 2); err != nil {
		return err
	}

	if err := validate.MaxLength("q", "query", o.Q, 128); err != nil {
		return err
	}

	return nil
}

// bindType binds and validates parameter Type from query.
func (o *SearchTorrentsParams) bindType(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false

	if raw == "" { // empty values pass all other validations
		return nil
	}
	o.Type = &raw

	if err := o.validateType(formats); err != nil {
		return err
	}

	return nil
}

// validateType carries on validations for parameter Type
func (o *SearchTorrentsParams) validateType(formats strfmt.Registry) error {

	if err := validate.EnumCase("type", "query", *o.Type, []interface{}{"movies", "music", "books", "others"}, true); err != nil {
		return err
	}

	return nil
}