// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/runtime/security"
	"github.com/go-openapi/spec"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"

	"github.com/RacoonMediaServer/rms-media-discovery/internal/server/models"
	"github.com/RacoonMediaServer/rms-media-discovery/internal/server/restapi/operations/accounts"
	"github.com/RacoonMediaServer/rms-media-discovery/internal/server/restapi/operations/movies"
	"github.com/RacoonMediaServer/rms-media-discovery/internal/server/restapi/operations/music"
	"github.com/RacoonMediaServer/rms-media-discovery/internal/server/restapi/operations/torrents"
)

// NewServerAPI creates a new Server instance
func NewServerAPI(spec *loads.Document) *ServerAPI {
	return &ServerAPI{
		handlers:            make(map[string]map[string]http.Handler),
		formats:             strfmt.Default,
		defaultConsumes:     "application/json",
		defaultProduces:     "application/json",
		customConsumers:     make(map[string]runtime.Consumer),
		customProducers:     make(map[string]runtime.Producer),
		PreServerShutdown:   func() {},
		ServerShutdown:      func() {},
		spec:                spec,
		useSwaggerUI:        false,
		ServeError:          errors.ServeError,
		BasicAuthenticator:  security.BasicAuth,
		APIKeyAuthenticator: security.APIKeyAuth,
		BearerAuthenticator: security.BearerAuth,

		JSONConsumer: runtime.JSONConsumer(),

		BinProducer:  runtime.ByteStreamProducer(),
		JSONProducer: runtime.JSONProducer(),

		AccountsCreateAccountHandler: accounts.CreateAccountHandlerFunc(func(params accounts.CreateAccountParams, principal *models.Principal) middleware.Responder {
			return middleware.NotImplemented("operation accounts.CreateAccount has not yet been implemented")
		}),
		AccountsDeleteAccountHandler: accounts.DeleteAccountHandlerFunc(func(params accounts.DeleteAccountParams, principal *models.Principal) middleware.Responder {
			return middleware.NotImplemented("operation accounts.DeleteAccount has not yet been implemented")
		}),
		TorrentsDownloadTorrentHandler: torrents.DownloadTorrentHandlerFunc(func(params torrents.DownloadTorrentParams, principal *models.Principal) middleware.Responder {
			return middleware.NotImplemented("operation torrents.DownloadTorrent has not yet been implemented")
		}),
		AccountsGetAccountsHandler: accounts.GetAccountsHandlerFunc(func(params accounts.GetAccountsParams, principal *models.Principal) middleware.Responder {
			return middleware.NotImplemented("operation accounts.GetAccounts has not yet been implemented")
		}),
		MoviesGetMovieInfoHandler: movies.GetMovieInfoHandlerFunc(func(params movies.GetMovieInfoParams, principal *models.Principal) middleware.Responder {
			return middleware.NotImplemented("operation movies.GetMovieInfo has not yet been implemented")
		}),
		MoviesSearchMoviesHandler: movies.SearchMoviesHandlerFunc(func(params movies.SearchMoviesParams, principal *models.Principal) middleware.Responder {
			return middleware.NotImplemented("operation movies.SearchMovies has not yet been implemented")
		}),
		MusicSearchMusicHandler: music.SearchMusicHandlerFunc(func(params music.SearchMusicParams, principal *models.Principal) middleware.Responder {
			return middleware.NotImplemented("operation music.SearchMusic has not yet been implemented")
		}),
		TorrentsSearchTorrentsHandler: torrents.SearchTorrentsHandlerFunc(func(params torrents.SearchTorrentsParams, principal *models.Principal) middleware.Responder {
			return middleware.NotImplemented("operation torrents.SearchTorrents has not yet been implemented")
		}),
		TorrentsSearchTorrentsAsyncHandler: torrents.SearchTorrentsAsyncHandlerFunc(func(params torrents.SearchTorrentsAsyncParams, principal *models.Principal) middleware.Responder {
			return middleware.NotImplemented("operation torrents.SearchTorrentsAsync has not yet been implemented")
		}),
		TorrentsSearchTorrentsAsyncCancelHandler: torrents.SearchTorrentsAsyncCancelHandlerFunc(func(params torrents.SearchTorrentsAsyncCancelParams, principal *models.Principal) middleware.Responder {
			return middleware.NotImplemented("operation torrents.SearchTorrentsAsyncCancel has not yet been implemented")
		}),
		TorrentsSearchTorrentsAsyncStatusHandler: torrents.SearchTorrentsAsyncStatusHandlerFunc(func(params torrents.SearchTorrentsAsyncStatusParams, principal *models.Principal) middleware.Responder {
			return middleware.NotImplemented("operation torrents.SearchTorrentsAsyncStatus has not yet been implemented")
		}),

		// Applies when the "x-token" header is set
		KeyAuth: func(token string) (*models.Principal, error) {
			return nil, errors.NotImplemented("api key auth (key) x-token from header param [x-token] has not yet been implemented")
		},
		// default authorizer is authorized meaning no requests are blocked
		APIAuthorizer: security.Authorized(),
	}
}

/*ServerAPI API for Racoon Media Server Project */
type ServerAPI struct {
	spec            *loads.Document
	context         *middleware.Context
	handlers        map[string]map[string]http.Handler
	formats         strfmt.Registry
	customConsumers map[string]runtime.Consumer
	customProducers map[string]runtime.Producer
	defaultConsumes string
	defaultProduces string
	Middleware      func(middleware.Builder) http.Handler
	useSwaggerUI    bool

	// BasicAuthenticator generates a runtime.Authenticator from the supplied basic auth function.
	// It has a default implementation in the security package, however you can replace it for your particular usage.
	BasicAuthenticator func(security.UserPassAuthentication) runtime.Authenticator

	// APIKeyAuthenticator generates a runtime.Authenticator from the supplied token auth function.
	// It has a default implementation in the security package, however you can replace it for your particular usage.
	APIKeyAuthenticator func(string, string, security.TokenAuthentication) runtime.Authenticator

	// BearerAuthenticator generates a runtime.Authenticator from the supplied bearer token auth function.
	// It has a default implementation in the security package, however you can replace it for your particular usage.
	BearerAuthenticator func(string, security.ScopedTokenAuthentication) runtime.Authenticator

	// JSONConsumer registers a consumer for the following mime types:
	//   - application/json
	JSONConsumer runtime.Consumer

	// BinProducer registers a producer for the following mime types:
	//   - application/octet-stream
	BinProducer runtime.Producer
	// JSONProducer registers a producer for the following mime types:
	//   - application/json
	JSONProducer runtime.Producer

	// KeyAuth registers a function that takes a token and returns a principal
	// it performs authentication based on an api key x-token provided in the header
	KeyAuth func(string) (*models.Principal, error)

	// APIAuthorizer provides access control (ACL/RBAC/ABAC) by providing access to the request and authenticated principal
	APIAuthorizer runtime.Authorizer

	// AccountsCreateAccountHandler sets the operation handler for the create account operation
	AccountsCreateAccountHandler accounts.CreateAccountHandler
	// AccountsDeleteAccountHandler sets the operation handler for the delete account operation
	AccountsDeleteAccountHandler accounts.DeleteAccountHandler
	// TorrentsDownloadTorrentHandler sets the operation handler for the download torrent operation
	TorrentsDownloadTorrentHandler torrents.DownloadTorrentHandler
	// AccountsGetAccountsHandler sets the operation handler for the get accounts operation
	AccountsGetAccountsHandler accounts.GetAccountsHandler
	// MoviesGetMovieInfoHandler sets the operation handler for the get movie info operation
	MoviesGetMovieInfoHandler movies.GetMovieInfoHandler
	// MoviesSearchMoviesHandler sets the operation handler for the search movies operation
	MoviesSearchMoviesHandler movies.SearchMoviesHandler
	// MusicSearchMusicHandler sets the operation handler for the search music operation
	MusicSearchMusicHandler music.SearchMusicHandler
	// TorrentsSearchTorrentsHandler sets the operation handler for the search torrents operation
	TorrentsSearchTorrentsHandler torrents.SearchTorrentsHandler
	// TorrentsSearchTorrentsAsyncHandler sets the operation handler for the search torrents async operation
	TorrentsSearchTorrentsAsyncHandler torrents.SearchTorrentsAsyncHandler
	// TorrentsSearchTorrentsAsyncCancelHandler sets the operation handler for the search torrents async cancel operation
	TorrentsSearchTorrentsAsyncCancelHandler torrents.SearchTorrentsAsyncCancelHandler
	// TorrentsSearchTorrentsAsyncStatusHandler sets the operation handler for the search torrents async status operation
	TorrentsSearchTorrentsAsyncStatusHandler torrents.SearchTorrentsAsyncStatusHandler

	// ServeError is called when an error is received, there is a default handler
	// but you can set your own with this
	ServeError func(http.ResponseWriter, *http.Request, error)

	// PreServerShutdown is called before the HTTP(S) server is shutdown
	// This allows for custom functions to get executed before the HTTP(S) server stops accepting traffic
	PreServerShutdown func()

	// ServerShutdown is called when the HTTP(S) server is shut down and done
	// handling all active connections and does not accept connections any more
	ServerShutdown func()

	// Custom command line argument groups with their descriptions
	CommandLineOptionsGroups []swag.CommandLineOptionsGroup

	// User defined logger function.
	Logger func(string, ...interface{})
}

// UseRedoc for documentation at /docs
func (o *ServerAPI) UseRedoc() {
	o.useSwaggerUI = false
}

// UseSwaggerUI for documentation at /docs
func (o *ServerAPI) UseSwaggerUI() {
	o.useSwaggerUI = true
}

// SetDefaultProduces sets the default produces media type
func (o *ServerAPI) SetDefaultProduces(mediaType string) {
	o.defaultProduces = mediaType
}

// SetDefaultConsumes returns the default consumes media type
func (o *ServerAPI) SetDefaultConsumes(mediaType string) {
	o.defaultConsumes = mediaType
}

// SetSpec sets a spec that will be served for the clients.
func (o *ServerAPI) SetSpec(spec *loads.Document) {
	o.spec = spec
}

// DefaultProduces returns the default produces media type
func (o *ServerAPI) DefaultProduces() string {
	return o.defaultProduces
}

// DefaultConsumes returns the default consumes media type
func (o *ServerAPI) DefaultConsumes() string {
	return o.defaultConsumes
}

// Formats returns the registered string formats
func (o *ServerAPI) Formats() strfmt.Registry {
	return o.formats
}

// RegisterFormat registers a custom format validator
func (o *ServerAPI) RegisterFormat(name string, format strfmt.Format, validator strfmt.Validator) {
	o.formats.Add(name, format, validator)
}

// Validate validates the registrations in the ServerAPI
func (o *ServerAPI) Validate() error {
	var unregistered []string

	if o.JSONConsumer == nil {
		unregistered = append(unregistered, "JSONConsumer")
	}

	if o.BinProducer == nil {
		unregistered = append(unregistered, "BinProducer")
	}
	if o.JSONProducer == nil {
		unregistered = append(unregistered, "JSONProducer")
	}

	if o.KeyAuth == nil {
		unregistered = append(unregistered, "XTokenAuth")
	}

	if o.AccountsCreateAccountHandler == nil {
		unregistered = append(unregistered, "accounts.CreateAccountHandler")
	}
	if o.AccountsDeleteAccountHandler == nil {
		unregistered = append(unregistered, "accounts.DeleteAccountHandler")
	}
	if o.TorrentsDownloadTorrentHandler == nil {
		unregistered = append(unregistered, "torrents.DownloadTorrentHandler")
	}
	if o.AccountsGetAccountsHandler == nil {
		unregistered = append(unregistered, "accounts.GetAccountsHandler")
	}
	if o.MoviesGetMovieInfoHandler == nil {
		unregistered = append(unregistered, "movies.GetMovieInfoHandler")
	}
	if o.MoviesSearchMoviesHandler == nil {
		unregistered = append(unregistered, "movies.SearchMoviesHandler")
	}
	if o.MusicSearchMusicHandler == nil {
		unregistered = append(unregistered, "music.SearchMusicHandler")
	}
	if o.TorrentsSearchTorrentsHandler == nil {
		unregistered = append(unregistered, "torrents.SearchTorrentsHandler")
	}
	if o.TorrentsSearchTorrentsAsyncHandler == nil {
		unregistered = append(unregistered, "torrents.SearchTorrentsAsyncHandler")
	}
	if o.TorrentsSearchTorrentsAsyncCancelHandler == nil {
		unregistered = append(unregistered, "torrents.SearchTorrentsAsyncCancelHandler")
	}
	if o.TorrentsSearchTorrentsAsyncStatusHandler == nil {
		unregistered = append(unregistered, "torrents.SearchTorrentsAsyncStatusHandler")
	}

	if len(unregistered) > 0 {
		return fmt.Errorf("missing registration: %s", strings.Join(unregistered, ", "))
	}

	return nil
}

// ServeErrorFor gets a error handler for a given operation id
func (o *ServerAPI) ServeErrorFor(operationID string) func(http.ResponseWriter, *http.Request, error) {
	return o.ServeError
}

// AuthenticatorsFor gets the authenticators for the specified security schemes
func (o *ServerAPI) AuthenticatorsFor(schemes map[string]spec.SecurityScheme) map[string]runtime.Authenticator {
	result := make(map[string]runtime.Authenticator)
	for name := range schemes {
		switch name {
		case "key":
			scheme := schemes[name]
			result[name] = o.APIKeyAuthenticator(scheme.Name, scheme.In, func(token string) (interface{}, error) {
				return o.KeyAuth(token)
			})

		}
	}
	return result
}

// Authorizer returns the registered authorizer
func (o *ServerAPI) Authorizer() runtime.Authorizer {
	return o.APIAuthorizer
}

// ConsumersFor gets the consumers for the specified media types.
// MIME type parameters are ignored here.
func (o *ServerAPI) ConsumersFor(mediaTypes []string) map[string]runtime.Consumer {
	result := make(map[string]runtime.Consumer, len(mediaTypes))
	for _, mt := range mediaTypes {
		switch mt {
		case "application/json":
			result["application/json"] = o.JSONConsumer
		}

		if c, ok := o.customConsumers[mt]; ok {
			result[mt] = c
		}
	}
	return result
}

// ProducersFor gets the producers for the specified media types.
// MIME type parameters are ignored here.
func (o *ServerAPI) ProducersFor(mediaTypes []string) map[string]runtime.Producer {
	result := make(map[string]runtime.Producer, len(mediaTypes))
	for _, mt := range mediaTypes {
		switch mt {
		case "application/octet-stream":
			result["application/octet-stream"] = o.BinProducer
		case "application/json":
			result["application/json"] = o.JSONProducer
		}

		if p, ok := o.customProducers[mt]; ok {
			result[mt] = p
		}
	}
	return result
}

// HandlerFor gets a http.Handler for the provided operation method and path
func (o *ServerAPI) HandlerFor(method, path string) (http.Handler, bool) {
	if o.handlers == nil {
		return nil, false
	}
	um := strings.ToUpper(method)
	if _, ok := o.handlers[um]; !ok {
		return nil, false
	}
	if path == "/" {
		path = ""
	}
	h, ok := o.handlers[um][path]
	return h, ok
}

// Context returns the middleware context for the server API
func (o *ServerAPI) Context() *middleware.Context {
	if o.context == nil {
		o.context = middleware.NewRoutableContext(o.spec, o, nil)
	}

	return o.context
}

func (o *ServerAPI) initHandlerCache() {
	o.Context() // don't care about the result, just that the initialization happened
	if o.handlers == nil {
		o.handlers = make(map[string]map[string]http.Handler)
	}

	if o.handlers["POST"] == nil {
		o.handlers["POST"] = make(map[string]http.Handler)
	}
	o.handlers["POST"]["/accounts"] = accounts.NewCreateAccount(o.context, o.AccountsCreateAccountHandler)
	if o.handlers["DELETE"] == nil {
		o.handlers["DELETE"] = make(map[string]http.Handler)
	}
	o.handlers["DELETE"]["/accounts/{id}"] = accounts.NewDeleteAccount(o.context, o.AccountsDeleteAccountHandler)
	if o.handlers["GET"] == nil {
		o.handlers["GET"] = make(map[string]http.Handler)
	}
	o.handlers["GET"]["/torrents/download"] = torrents.NewDownloadTorrent(o.context, o.TorrentsDownloadTorrentHandler)
	if o.handlers["GET"] == nil {
		o.handlers["GET"] = make(map[string]http.Handler)
	}
	o.handlers["GET"]["/accounts"] = accounts.NewGetAccounts(o.context, o.AccountsGetAccountsHandler)
	if o.handlers["GET"] == nil {
		o.handlers["GET"] = make(map[string]http.Handler)
	}
	o.handlers["GET"]["/movies/{id}"] = movies.NewGetMovieInfo(o.context, o.MoviesGetMovieInfoHandler)
	if o.handlers["GET"] == nil {
		o.handlers["GET"] = make(map[string]http.Handler)
	}
	o.handlers["GET"]["/movies/search"] = movies.NewSearchMovies(o.context, o.MoviesSearchMoviesHandler)
	if o.handlers["GET"] == nil {
		o.handlers["GET"] = make(map[string]http.Handler)
	}
	o.handlers["GET"]["/music/search"] = music.NewSearchMusic(o.context, o.MusicSearchMusicHandler)
	if o.handlers["GET"] == nil {
		o.handlers["GET"] = make(map[string]http.Handler)
	}
	o.handlers["GET"]["/torrents/search"] = torrents.NewSearchTorrents(o.context, o.TorrentsSearchTorrentsHandler)
	if o.handlers["POST"] == nil {
		o.handlers["POST"] = make(map[string]http.Handler)
	}
	o.handlers["POST"]["/torrents/search:run"] = torrents.NewSearchTorrentsAsync(o.context, o.TorrentsSearchTorrentsAsyncHandler)
	if o.handlers["POST"] == nil {
		o.handlers["POST"] = make(map[string]http.Handler)
	}
	o.handlers["POST"]["/torrents/search/{id}:cancel"] = torrents.NewSearchTorrentsAsyncCancel(o.context, o.TorrentsSearchTorrentsAsyncCancelHandler)
	if o.handlers["GET"] == nil {
		o.handlers["GET"] = make(map[string]http.Handler)
	}
	o.handlers["GET"]["/torrents/search/{id}:status"] = torrents.NewSearchTorrentsAsyncStatus(o.context, o.TorrentsSearchTorrentsAsyncStatusHandler)
}

// Serve creates a http handler to serve the API over HTTP
// can be used directly in http.ListenAndServe(":8000", api.Serve(nil))
func (o *ServerAPI) Serve(builder middleware.Builder) http.Handler {
	o.Init()

	if o.Middleware != nil {
		return o.Middleware(builder)
	}
	if o.useSwaggerUI {
		return o.context.APIHandlerSwaggerUI(builder)
	}
	return o.context.APIHandler(builder)
}

// Init allows you to just initialize the handler cache, you can then recompose the middleware as you see fit
func (o *ServerAPI) Init() {
	if len(o.handlers) == 0 {
		o.initHandlerCache()
	}
}

// RegisterConsumer allows you to add (or override) a consumer for a media type.
func (o *ServerAPI) RegisterConsumer(mediaType string, consumer runtime.Consumer) {
	o.customConsumers[mediaType] = consumer
}

// RegisterProducer allows you to add (or override) a producer for a media type.
func (o *ServerAPI) RegisterProducer(mediaType string, producer runtime.Producer) {
	o.customProducers[mediaType] = producer
}

// AddMiddlewareFor adds a http middleware to existing handler
func (o *ServerAPI) AddMiddlewareFor(method, path string, builder middleware.Builder) {
	um := strings.ToUpper(method)
	if path == "/" {
		path = ""
	}
	o.Init()
	if h, ok := o.handlers[um][path]; ok {
		o.handlers[um][path] = builder(h)
	}
}
