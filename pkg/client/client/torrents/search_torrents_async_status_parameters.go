// Code generated by go-swagger; DO NOT EDIT.

package torrents

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
)

// NewSearchTorrentsAsyncStatusParams creates a new SearchTorrentsAsyncStatusParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewSearchTorrentsAsyncStatusParams() *SearchTorrentsAsyncStatusParams {
	return &SearchTorrentsAsyncStatusParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewSearchTorrentsAsyncStatusParamsWithTimeout creates a new SearchTorrentsAsyncStatusParams object
// with the ability to set a timeout on a request.
func NewSearchTorrentsAsyncStatusParamsWithTimeout(timeout time.Duration) *SearchTorrentsAsyncStatusParams {
	return &SearchTorrentsAsyncStatusParams{
		timeout: timeout,
	}
}

// NewSearchTorrentsAsyncStatusParamsWithContext creates a new SearchTorrentsAsyncStatusParams object
// with the ability to set a context for a request.
func NewSearchTorrentsAsyncStatusParamsWithContext(ctx context.Context) *SearchTorrentsAsyncStatusParams {
	return &SearchTorrentsAsyncStatusParams{
		Context: ctx,
	}
}

// NewSearchTorrentsAsyncStatusParamsWithHTTPClient creates a new SearchTorrentsAsyncStatusParams object
// with the ability to set a custom HTTPClient for a request.
func NewSearchTorrentsAsyncStatusParamsWithHTTPClient(client *http.Client) *SearchTorrentsAsyncStatusParams {
	return &SearchTorrentsAsyncStatusParams{
		HTTPClient: client,
	}
}

/*
SearchTorrentsAsyncStatusParams contains all the parameters to send to the API endpoint

	for the search torrents async status operation.

	Typically these are written to a http.Request.
*/
type SearchTorrentsAsyncStatusParams struct {

	// ID.
	ID string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the search torrents async status params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *SearchTorrentsAsyncStatusParams) WithDefaults() *SearchTorrentsAsyncStatusParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the search torrents async status params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *SearchTorrentsAsyncStatusParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the search torrents async status params
func (o *SearchTorrentsAsyncStatusParams) WithTimeout(timeout time.Duration) *SearchTorrentsAsyncStatusParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the search torrents async status params
func (o *SearchTorrentsAsyncStatusParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the search torrents async status params
func (o *SearchTorrentsAsyncStatusParams) WithContext(ctx context.Context) *SearchTorrentsAsyncStatusParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the search torrents async status params
func (o *SearchTorrentsAsyncStatusParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the search torrents async status params
func (o *SearchTorrentsAsyncStatusParams) WithHTTPClient(client *http.Client) *SearchTorrentsAsyncStatusParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the search torrents async status params
func (o *SearchTorrentsAsyncStatusParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithID adds the id to the search torrents async status params
func (o *SearchTorrentsAsyncStatusParams) WithID(id string) *SearchTorrentsAsyncStatusParams {
	o.SetID(id)
	return o
}

// SetID adds the id to the search torrents async status params
func (o *SearchTorrentsAsyncStatusParams) SetID(id string) {
	o.ID = id
}

// WriteToRequest writes these params to a swagger request
func (o *SearchTorrentsAsyncStatusParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param id
	if err := r.SetPathParam("id", o.ID); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}