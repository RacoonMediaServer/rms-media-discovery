// Code generated by go-swagger; DO NOT EDIT.

package music

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// NewSearchMusicParams creates a new SearchMusicParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewSearchMusicParams() *SearchMusicParams {
	return &SearchMusicParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewSearchMusicParamsWithTimeout creates a new SearchMusicParams object
// with the ability to set a timeout on a request.
func NewSearchMusicParamsWithTimeout(timeout time.Duration) *SearchMusicParams {
	return &SearchMusicParams{
		timeout: timeout,
	}
}

// NewSearchMusicParamsWithContext creates a new SearchMusicParams object
// with the ability to set a context for a request.
func NewSearchMusicParamsWithContext(ctx context.Context) *SearchMusicParams {
	return &SearchMusicParams{
		Context: ctx,
	}
}

// NewSearchMusicParamsWithHTTPClient creates a new SearchMusicParams object
// with the ability to set a custom HTTPClient for a request.
func NewSearchMusicParamsWithHTTPClient(client *http.Client) *SearchMusicParams {
	return &SearchMusicParams{
		HTTPClient: client,
	}
}

/*
SearchMusicParams contains all the parameters to send to the API endpoint

	for the search music operation.

	Typically these are written to a http.Request.
*/
type SearchMusicParams struct {

	/* Limit.

	   Ограничение на кол-во результатов
	*/
	Limit *int64

	/* Q.

	   Искомый запрос
	*/
	Q string

	// Type.
	//
	// Default: "any"
	Type *string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the search music params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *SearchMusicParams) WithDefaults() *SearchMusicParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the search music params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *SearchMusicParams) SetDefaults() {
	var (
		typeVarDefault = string("any")
	)

	val := SearchMusicParams{
		Type: &typeVarDefault,
	}

	val.timeout = o.timeout
	val.Context = o.Context
	val.HTTPClient = o.HTTPClient
	*o = val
}

// WithTimeout adds the timeout to the search music params
func (o *SearchMusicParams) WithTimeout(timeout time.Duration) *SearchMusicParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the search music params
func (o *SearchMusicParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the search music params
func (o *SearchMusicParams) WithContext(ctx context.Context) *SearchMusicParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the search music params
func (o *SearchMusicParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the search music params
func (o *SearchMusicParams) WithHTTPClient(client *http.Client) *SearchMusicParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the search music params
func (o *SearchMusicParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithLimit adds the limit to the search music params
func (o *SearchMusicParams) WithLimit(limit *int64) *SearchMusicParams {
	o.SetLimit(limit)
	return o
}

// SetLimit adds the limit to the search music params
func (o *SearchMusicParams) SetLimit(limit *int64) {
	o.Limit = limit
}

// WithQ adds the q to the search music params
func (o *SearchMusicParams) WithQ(q string) *SearchMusicParams {
	o.SetQ(q)
	return o
}

// SetQ adds the q to the search music params
func (o *SearchMusicParams) SetQ(q string) {
	o.Q = q
}

// WithType adds the typeVar to the search music params
func (o *SearchMusicParams) WithType(typeVar *string) *SearchMusicParams {
	o.SetType(typeVar)
	return o
}

// SetType adds the type to the search music params
func (o *SearchMusicParams) SetType(typeVar *string) {
	o.Type = typeVar
}

// WriteToRequest writes these params to a swagger request
func (o *SearchMusicParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.Limit != nil {

		// query param limit
		var qrLimit int64

		if o.Limit != nil {
			qrLimit = *o.Limit
		}
		qLimit := swag.FormatInt64(qrLimit)
		if qLimit != "" {

			if err := r.SetQueryParam("limit", qLimit); err != nil {
				return err
			}
		}
	}

	// query param q
	qrQ := o.Q
	qQ := qrQ
	if qQ != "" {

		if err := r.SetQueryParam("q", qQ); err != nil {
			return err
		}
	}

	if o.Type != nil {

		// query param type
		var qrType string

		if o.Type != nil {
			qrType = *o.Type
		}
		qType := qrType
		if qType != "" {

			if err := r.SetQueryParam("type", qType); err != nil {
				return err
			}
		}
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
