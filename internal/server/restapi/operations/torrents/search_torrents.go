// Code generated by go-swagger; DO NOT EDIT.

package torrents

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"

	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/server/models"
)

// SearchTorrentsHandlerFunc turns a function with the right signature into a search torrents handler
type SearchTorrentsHandlerFunc func(SearchTorrentsParams, *models.Principal) middleware.Responder

// Handle executing the request and returning a response
func (fn SearchTorrentsHandlerFunc) Handle(params SearchTorrentsParams, principal *models.Principal) middleware.Responder {
	return fn(params, principal)
}

// SearchTorrentsHandler interface for that can handle valid search torrents params
type SearchTorrentsHandler interface {
	Handle(SearchTorrentsParams, *models.Principal) middleware.Responder
}

// NewSearchTorrents creates a new http.Handler for the search torrents operation
func NewSearchTorrents(ctx *middleware.Context, handler SearchTorrentsHandler) *SearchTorrents {
	return &SearchTorrents{Context: ctx, Handler: handler}
}

/*
	SearchTorrents swagger:route GET /torrents/search torrents searchTorrents

# Поиск контента на торрент-трекерах

Поиск фильмов и сериалов по названию на различных платформах
*/
type SearchTorrents struct {
	Context *middleware.Context
	Handler SearchTorrentsHandler
}

func (o *SearchTorrents) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewSearchTorrentsParams()
	uprinc, aCtx, err := o.Context.Authorize(r, route)
	if err != nil {
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}
	if aCtx != nil {
		*r = *aCtx
	}
	var principal *models.Principal
	if uprinc != nil {
		principal = uprinc.(*models.Principal) // this is really a models.Principal, I promise
	}

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params, principal) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}

// SearchTorrentsOKBody search torrents o k body
//
// swagger:model SearchTorrentsOKBody
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
