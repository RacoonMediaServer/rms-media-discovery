package server

import (
	"context"
	"net/http"

	"github.com/RacoonMediaServer/rms-media-discovery/internal/server/restapi/operations/music"
	rms_users "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-users"

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
	api.TorrentsSearchTorrentsAsyncHandler = torrents.SearchTorrentsAsyncHandlerFunc(s.searchTorrentsAsync)
	api.TorrentsSearchTorrentsAsyncStatusHandler = torrents.SearchTorrentsAsyncStatusHandlerFunc(s.searchTorrentsAsyncStatus)
	api.TorrentsSearchTorrentsAsyncCancelHandler = torrents.SearchTorrentsAsyncCancelHandlerFunc(s.searchTorrentsAsyncCancel)
	api.TorrentsDownloadTorrentHandler = torrents.DownloadTorrentHandlerFunc(s.downloadTorrent)

	api.AccountsGetAccountsHandler = accounts.GetAccountsHandlerFunc(s.getAccounts)
	api.AccountsCreateAccountHandler = accounts.CreateAccountHandlerFunc(s.createAccount)
	api.AccountsDeleteAccountHandler = accounts.DeleteAccountHandlerFunc(s.deleteAccount)

	api.KeyAuth = func(token string) (*models.Principal, error) {
		req := rms_users.CheckPermissionsRequest{
			Token: token,
			Perms: []rms_users.Permissions{rms_users.Permissions_Search},
		}
		resp, err := s.Users.CheckPermissions(context.Background(), &req)
		if err != nil {
			s.log.Errorf("Cannot retrieve permissions: %s", err)
			return nil, errors.New(http.StatusForbidden, "Forbidden")
		}
		if !resp.Allowed {
			return nil, errors.New(http.StatusForbidden, "Forbidden")
		}

		userID := resp.UserId
		req.Perms = []rms_users.Permissions{rms_users.Permissions_AccountManagement}
		resp, err = s.Users.CheckPermissions(context.Background(), &req)
		if err != nil {
			s.log.Errorf("Cannot retrieve permissions: %s", err)
			return nil, errors.New(http.StatusForbidden, "Forbidden")
		}

		return &models.Principal{Token: userID, CanManageAccounts: resp.Allowed}, nil
	}
}
