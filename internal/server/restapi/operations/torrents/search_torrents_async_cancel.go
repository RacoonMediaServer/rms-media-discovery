// Code generated by go-swagger; DO NOT EDIT.

package torrents

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"

	"github.com/RacoonMediaServer/rms-media-discovery/internal/server/models"
)

// SearchTorrentsAsyncCancelHandlerFunc turns a function with the right signature into a search torrents async cancel handler
type SearchTorrentsAsyncCancelHandlerFunc func(SearchTorrentsAsyncCancelParams, *models.Principal) middleware.Responder

// Handle executing the request and returning a response
func (fn SearchTorrentsAsyncCancelHandlerFunc) Handle(params SearchTorrentsAsyncCancelParams, principal *models.Principal) middleware.Responder {
	return fn(params, principal)
}

// SearchTorrentsAsyncCancelHandler interface for that can handle valid search torrents async cancel params
type SearchTorrentsAsyncCancelHandler interface {
	Handle(SearchTorrentsAsyncCancelParams, *models.Principal) middleware.Responder
}

// NewSearchTorrentsAsyncCancel creates a new http.Handler for the search torrents async cancel operation
func NewSearchTorrentsAsyncCancel(ctx *middleware.Context, handler SearchTorrentsAsyncCancelHandler) *SearchTorrentsAsyncCancel {
	return &SearchTorrentsAsyncCancel{Context: ctx, Handler: handler}
}

/*
	SearchTorrentsAsyncCancel swagger:route POST /torrents/search/{id}:cancel torrents searchTorrentsAsyncCancel

# Отменить задачу

Отмена и удаление задачи поиска
*/
type SearchTorrentsAsyncCancel struct {
	Context *middleware.Context
	Handler SearchTorrentsAsyncCancelHandler
}

func (o *SearchTorrentsAsyncCancel) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewSearchTorrentsAsyncCancelParams()
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
