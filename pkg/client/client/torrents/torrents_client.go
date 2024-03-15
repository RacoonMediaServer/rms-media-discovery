// Code generated by go-swagger; DO NOT EDIT.

package torrents

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
)

// New creates a new torrents API client.
func New(transport runtime.ClientTransport, formats strfmt.Registry) ClientService {
	return &Client{transport: transport, formats: formats}
}

/*
Client for torrents API
*/
type Client struct {
	transport runtime.ClientTransport
	formats   strfmt.Registry
}

// ClientOption is the option for Client methods
type ClientOption func(*runtime.ClientOperation)

// ClientService is the interface for Client methods
type ClientService interface {
	PostTorrentSearchIDCancel(params *PostTorrentSearchIDCancelParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*PostTorrentSearchIDCancelOK, error)

	DownloadTorrent(params *DownloadTorrentParams, authInfo runtime.ClientAuthInfoWriter, writer io.Writer, opts ...ClientOption) (*DownloadTorrentOK, error)

	GetSearchTorrentsStatus(params *GetSearchTorrentsStatusParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*GetSearchTorrentsStatusOK, error)

	SearchTorrents(params *SearchTorrentsParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*SearchTorrentsOK, error)

	SearchTorrentsAsync(params *SearchTorrentsAsyncParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*SearchTorrentsAsyncOK, error)

	SetTransport(transport runtime.ClientTransport)
}

/*
PostTorrentSearchIDCancel отменитьs задачу

Отмена и удаление задачи поиска
*/
func (a *Client) PostTorrentSearchIDCancel(params *PostTorrentSearchIDCancelParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*PostTorrentSearchIDCancelOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewPostTorrentSearchIDCancelParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "PostTorrentSearchIDCancel",
		Method:             "POST",
		PathPattern:        "/torrent/search/{id}:cancel",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &PostTorrentSearchIDCancelReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*PostTorrentSearchIDCancelOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for PostTorrentSearchIDCancel: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
DownloadTorrent загрузкаs торрент файла

Позволяет скачать торрент-файл, с помощью которого можно скачать контент
*/
func (a *Client) DownloadTorrent(params *DownloadTorrentParams, authInfo runtime.ClientAuthInfoWriter, writer io.Writer, opts ...ClientOption) (*DownloadTorrentOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewDownloadTorrentParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "downloadTorrent",
		Method:             "GET",
		PathPattern:        "/torrents/download",
		ProducesMediaTypes: []string{"application/octet-stream"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &DownloadTorrentReader{formats: a.formats, writer: writer},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*DownloadTorrentOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for downloadTorrent: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
GetSearchTorrentsStatus узнатьs статус задачи поиска

Запросить статус и результаты задачи поиска
*/
func (a *Client) GetSearchTorrentsStatus(params *GetSearchTorrentsStatusParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*GetSearchTorrentsStatusOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetSearchTorrentsStatusParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "getSearchTorrentsStatus",
		Method:             "GET",
		PathPattern:        "/torrents/search/{id}:status",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &GetSearchTorrentsStatusReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*GetSearchTorrentsStatusOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for getSearchTorrentsStatus: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
SearchTorrents поискs контента на торрент трекерах

Поиск раздач на различных платформах
*/
func (a *Client) SearchTorrents(params *SearchTorrentsParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*SearchTorrentsOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewSearchTorrentsParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "searchTorrents",
		Method:             "GET",
		PathPattern:        "/torrents/search",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &SearchTorrentsReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*SearchTorrentsOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for searchTorrents: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
SearchTorrentsAsync стартs задачи поиска раздач

LRO поиск раздач на торрент-трекерах
*/
func (a *Client) SearchTorrentsAsync(params *SearchTorrentsAsyncParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*SearchTorrentsAsyncOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewSearchTorrentsAsyncParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "searchTorrentsAsync",
		Method:             "POST",
		PathPattern:        "/torrents/search:run",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &SearchTorrentsAsyncReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*SearchTorrentsAsyncOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for searchTorrentsAsync: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

// SetTransport changes the transport on the client
func (a *Client) SetTransport(transport runtime.ClientTransport) {
	a.transport = transport
}
