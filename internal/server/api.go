package server

import (
	"context"
	"github.com/RacoonMediaServer/rms-media-discovery/internal/server/restapi/operations/music"
	rms_users "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-users"
	"net/http"

	"github.com/RacoonMediaServer/rms-media-discovery/internal/server/models"
	"github.com/RacoonMediaServer/rms-media-discovery/internal/server/restapi/operations"
	"github.com/RacoonMediaServer/rms-media-discovery/internal/server/restapi/operations/accounts"
	"github.com/RacoonMediaServer/rms-media-discovery/internal/server/restapi/operations/movies"
	"github.com/RacoonMediaServer/rms-media-discovery/internal/server/restapi/operations/torrents"
	"github.com/go-openapi/errors"
)

func (s *Server) configureAPI(api *operations.ServerAPI) {
	api.MoviesSearchMoviesHandler = movies.SearchMoviesHandlerFunc(s.searchMovies)
	api.MoviesGetMovieInfoHandler = movies.GetMovieInfoHandlerFunc(s.getMovieInfo)

	api.MusicSearchMusicHandler = music.SearchMusicHandlerFunc(s.searchMusic)

	api.TorrentsSearchTorrentsHandler = torrents.SearchTorrentsHandlerFunc(s.searchTorrents)
	api.TorrentsDownloadTorrentHandler = torrents.DownloadTorrentHandlerFunc(s.downloadTorrent)

	api.AccountsGetAccountsHandler = accounts.GetAccountsHandlerFunc(s.getAccounts)
	api.AccountsCreateAccountHandler = accounts.CreateAccountHandlerFunc(s.createAccount)
	api.AccountsDeleteAccountHandler = accounts.DeleteAccountHandlerFunc(s.deleteAccount)

	api.KeyAuth = func(token string) (*models.Principal, error) {
		resp, err := s.Users.GetPermissions(context.Background(), &rms_users.GetPermissionsRequest{Token: token})
		if err != nil {
			s.log.Errorf("Cannot retrieve permissions: %s", err)
			return nil, errors.New(http.StatusForbidden, "Forbidden")
		}
		searchAllowed := false
		manageAllowed := false
		for _, p := range resp.Perms {
			switch p {
			case rms_users.Permissions_Search:
				searchAllowed = true
			case rms_users.Permissions_AccountManagement:
				manageAllowed = true
			}
		}
		if !searchAllowed {
			return nil, errors.New(http.StatusForbidden, "Forbidden")
		}
		return &models.Principal{Token: token, CanManageAccounts: manageAllowed}, nil
	}
}
