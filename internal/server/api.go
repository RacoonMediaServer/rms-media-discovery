package server

import (
	"net/http"

	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/server/models"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/server/restapi/operations"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/server/restapi/operations/accounts"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/server/restapi/operations/movies"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/server/restapi/operations/torrents"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/server/restapi/operations/users"
	"github.com/go-openapi/errors"
)

func (s *Server) configureAPI(api *operations.ServerAPI) {
	api.MoviesSearchMoviesHandler = movies.SearchMoviesHandlerFunc(s.searchMovies)
	api.TorrentsSearchTorrentsHandler = torrents.SearchTorrentsHandlerFunc(s.searchTorrents)
	api.TorrentsDownloadTorrentHandler = torrents.DownloadTorrentHandlerFunc(s.downloadTorrent)

	api.UsersGetUsersHandler = users.GetUsersHandlerFunc(s.getUsers)
	api.UsersCreateUserHandler = users.CreateUserHandlerFunc(s.createUser)
	api.UsersDeleteUserHandler = users.DeleteUserHandlerFunc(s.deleteUser)

	api.AccountsGetAccountsHandler = accounts.GetAccountsHandlerFunc(s.getAccounts)
	api.AccountsCreateAccountHandler = accounts.CreateAccountHandlerFunc(s.createAccount)
	api.AccountsDeleteAccountHandler = accounts.DeleteAccountHandlerFunc(s.deleteAccount)

	api.KeyAuth = func(token string) (*models.Principal, error) {
		valid, admin := s.Users.CheckAccess(token)
		if !valid {
			return nil, errors.New(http.StatusForbidden, "Forbidden")
		}
		return &models.Principal{Token: token, Admin: admin}, nil
	}

}
