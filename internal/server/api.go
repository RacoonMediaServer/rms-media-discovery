package server

import (
	"net/http"

	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/server/models"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/server/restapi/operations"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/server/restapi/operations/admin"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/server/restapi/operations/movies"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/server/restapi/operations/torrents"
	"github.com/go-openapi/errors"
)

func (s *Server) configureAPI(api *operations.ServerAPI) {
	api.MoviesSearchMoviesHandler = movies.SearchMoviesHandlerFunc(s.searchMovies)
	api.TorrentsSearchTorrentsHandler = torrents.SearchTorrentsHandlerFunc(s.searchTorrents)
	api.TorrentsDownloadTorrentHandler = torrents.DownloadTorrentHandlerFunc(s.downloadTorrent)

	api.AdminGetUsersHandler = admin.GetUsersHandlerFunc(s.getUsers)
	api.AdminCreateUserHandler = admin.CreateUserHandlerFunc(s.createUser)
	api.AdminDeleteUserHandler = admin.DeleteUserHandlerFunc(s.deleteUser)

	api.AdminGetAccountsHandler = admin.GetAccountsHandlerFunc(s.getAccounts)
	api.AdminCreateAccountHandler = admin.CreateAccountHandlerFunc(s.createAccount)
	api.AdminDeleteAccountHandler = admin.DeleteAccountHandlerFunc(s.deleteAccount)

	api.KeyAuth = func(token string) (*models.Principal, error) {
		valid, admin := s.Admin.CheckAccess(token)
		if !valid {
			return nil, errors.New(http.StatusForbidden, "Forbidden")
		}
		return &models.Principal{Token: token, Admin: admin}, nil
	}

}
