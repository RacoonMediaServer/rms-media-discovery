package server

import (
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/server/models"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/server/restapi/operations/torrents"
	torrents2 "git.rms.local/RacoonMediaServer/rms-media-discovery/internal/service/torrents"
	"github.com/go-openapi/runtime/middleware"
)

func (s *Server) searchTorrents(params torrents.SearchTorrentsParams, key *models.Principal) middleware.Responder {
	var limit uint
	if params.Limit != nil {
		limit = uint(*params.Limit)
	}

	hint := torrents2.SearchType_Other
	if params.Type != nil {
		switch *params.Type {
		case "movies":
			hint = torrents2.SearchType_Movies
		case "music":
			hint = torrents2.SearchType_Music
		case "books":
			hint = torrents2.SearchType_Books
		}
	}
	mov, err := s.Torrents.Search(params.HTTPRequest.Context(), params.Q, hint, limit)
	if err != nil {
		s.log.WithField("query", params.Q).Errorf("Search failed: %w", err)
		return torrents.NewSearchTorrentsInternalServerError()
	}

	payload := torrents.SearchTorrentsOKBody{Results: []*models.SearchTorrentsResult{}}
	for i := range mov {
		// TODO: менеджмент ссылок и остальные параметры
		result := &models.SearchTorrentsResult{
			Link:    &mov[i].Link,
			Seeders: int64(mov[i].Seeders),
			Size:    0,
			Title:   mov[i].Title,
		}
		payload.Results = append(payload.Results, result)
	}

	return torrents.NewSearchTorrentsOK().WithPayload(&payload)
}

func (s *Server) downloadTorrent(params torrents.DownloadTorrentParams, key *models.Principal) middleware.Responder {
	return torrents.NewDownloadTorrentInternalServerError()
}
