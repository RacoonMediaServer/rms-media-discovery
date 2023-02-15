// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	"github.com/RacoonMediaServer/rms-media-discovery/internal/server/models"
	"github.com/RacoonMediaServer/rms-media-discovery/internal/server/restapi/operations"
	"github.com/RacoonMediaServer/rms-media-discovery/internal/server/restapi/operations/accounts"
	"github.com/RacoonMediaServer/rms-media-discovery/internal/server/restapi/operations/movies"
	"github.com/RacoonMediaServer/rms-media-discovery/internal/server/restapi/operations/torrents"
	"github.com/RacoonMediaServer/rms-media-discovery/internal/server/restapi/operations/users"
)

//go:generate swagger generate server --target ../../server --name Server --spec ../../../api/discovery.yml --principal models.Principal --exclude-main

func configureFlags(api *operations.ServerAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.ServerAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.UseSwaggerUI()
	// To continue using redoc as your UI, uncomment the following line
	// api.UseRedoc()

	api.JSONConsumer = runtime.JSONConsumer()

	api.BinProducer = runtime.ByteStreamProducer()
	api.JSONProducer = runtime.JSONProducer()

	// Applies when the "x-token" header is set
	if api.KeyAuth == nil {
		api.KeyAuth = func(token string) (*models.Principal, error) {
			return nil, errors.NotImplemented("api key auth (key) x-token from header param [x-token] has not yet been implemented")
		}
	}

	// Set your custom authorizer if needed. Default one is security.Authorized()
	// Expected interface runtime.Authorizer
	//
	// Example:
	// api.APIAuthorizer = security.Authorized()

	if api.AccountsCreateAccountHandler == nil {
		api.AccountsCreateAccountHandler = accounts.CreateAccountHandlerFunc(func(params accounts.CreateAccountParams, principal *models.Principal) middleware.Responder {
			return middleware.NotImplemented("operation accounts.CreateAccount has not yet been implemented")
		})
	}
	if api.UsersCreateUserHandler == nil {
		api.UsersCreateUserHandler = users.CreateUserHandlerFunc(func(params users.CreateUserParams, principal *models.Principal) middleware.Responder {
			return middleware.NotImplemented("operation users.CreateUser has not yet been implemented")
		})
	}
	if api.AccountsDeleteAccountHandler == nil {
		api.AccountsDeleteAccountHandler = accounts.DeleteAccountHandlerFunc(func(params accounts.DeleteAccountParams, principal *models.Principal) middleware.Responder {
			return middleware.NotImplemented("operation accounts.DeleteAccount has not yet been implemented")
		})
	}
	if api.UsersDeleteUserHandler == nil {
		api.UsersDeleteUserHandler = users.DeleteUserHandlerFunc(func(params users.DeleteUserParams, principal *models.Principal) middleware.Responder {
			return middleware.NotImplemented("operation users.DeleteUser has not yet been implemented")
		})
	}
	if api.TorrentsDownloadTorrentHandler == nil {
		api.TorrentsDownloadTorrentHandler = torrents.DownloadTorrentHandlerFunc(func(params torrents.DownloadTorrentParams, principal *models.Principal) middleware.Responder {
			return middleware.NotImplemented("operation torrents.DownloadTorrent has not yet been implemented")
		})
	}
	if api.AccountsGetAccountsHandler == nil {
		api.AccountsGetAccountsHandler = accounts.GetAccountsHandlerFunc(func(params accounts.GetAccountsParams, principal *models.Principal) middleware.Responder {
			return middleware.NotImplemented("operation accounts.GetAccounts has not yet been implemented")
		})
	}
	if api.UsersGetUsersHandler == nil {
		api.UsersGetUsersHandler = users.GetUsersHandlerFunc(func(params users.GetUsersParams, principal *models.Principal) middleware.Responder {
			return middleware.NotImplemented("operation users.GetUsers has not yet been implemented")
		})
	}
	if api.MoviesSearchMoviesHandler == nil {
		api.MoviesSearchMoviesHandler = movies.SearchMoviesHandlerFunc(func(params movies.SearchMoviesParams, principal *models.Principal) middleware.Responder {
			return middleware.NotImplemented("operation movies.SearchMovies has not yet been implemented")
		})
	}
	if api.TorrentsSearchTorrentsHandler == nil {
		api.TorrentsSearchTorrentsHandler = torrents.SearchTorrentsHandlerFunc(func(params torrents.SearchTorrentsParams, principal *models.Principal) middleware.Responder {
			return middleware.NotImplemented("operation torrents.SearchTorrents has not yet been implemented")
		})
	}

	api.PreServerShutdown = func() {}

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix".
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation.
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics.
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
