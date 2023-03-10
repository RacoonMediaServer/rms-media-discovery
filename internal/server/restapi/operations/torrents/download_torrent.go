// Code generated by go-swagger; DO NOT EDIT.

package torrents

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"

	"github.com/RacoonMediaServer/rms-media-discovery/internal/server/models"
)

// DownloadTorrentHandlerFunc turns a function with the right signature into a download torrent handler
type DownloadTorrentHandlerFunc func(DownloadTorrentParams, *models.Principal) middleware.Responder

// Handle executing the request and returning a response
func (fn DownloadTorrentHandlerFunc) Handle(params DownloadTorrentParams, principal *models.Principal) middleware.Responder {
	return fn(params, principal)
}

// DownloadTorrentHandler interface for that can handle valid download torrent params
type DownloadTorrentHandler interface {
	Handle(DownloadTorrentParams, *models.Principal) middleware.Responder
}

// NewDownloadTorrent creates a new http.Handler for the download torrent operation
func NewDownloadTorrent(ctx *middleware.Context, handler DownloadTorrentHandler) *DownloadTorrent {
	return &DownloadTorrent{Context: ctx, Handler: handler}
}

/*
	DownloadTorrent swagger:route GET /torrents/download torrents downloadTorrent

# Загрузка торрент-файла

Позволяет скачать торрент-файл, с помощью которого можно скачать контент
*/
type DownloadTorrent struct {
	Context *middleware.Context
	Handler DownloadTorrentHandler
}

func (o *DownloadTorrent) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewDownloadTorrentParams()
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
