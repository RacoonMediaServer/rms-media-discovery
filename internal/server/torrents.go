package server

import (
	"bytes"
	"errors"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/server/models"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/server/restapi/operations/torrents"
	torrents2 "git.rms.local/RacoonMediaServer/rms-media-discovery/internal/service/torrents"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/pkg/media"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/pkg/model"
	"github.com/go-openapi/runtime/middleware"
	"io"
)

func convertTorrent(t *model.Torrent, contentType media.ContentType) *models.SearchTorrentsResult {
	seeders := int64(t.Seeders)
	size := int64(t.SizeMB)

	result := &models.SearchTorrentsResult{
		Link:    &t.Link,
		Seeders: &seeders,
		Size:    &size,
		Title:   &t.Title,
	}

	result.Format = t.Info.Format

	if contentType == media.Movies {
		for _, s := range t.Info.Seasons {
			result.Seasons = append(result.Seasons, int64(s))
		}
		result.Rip = t.Info.Rip
		result.Voice = t.Info.Voice
		result.Subtitles = t.Info.Subtitles
		result.Quality = t.Info.Quality.String()
	}

	return result
}

func searchQueryFromParams(params *torrents.SearchTorrentsParams) model.SearchQuery {
	var limit uint
	if params.Limit != nil {
		limit = uint(*params.Limit)
	}
	var strong bool
	if params.Strong != nil {
		strong = *params.Strong
	}

	hint := media.Other
	if params.Type != nil {
		switch *params.Type {
		case "movies":
			hint = media.Movies
		case "music":
			hint = media.Music
		case "books":
			hint = media.Books
		}
	}

	var year uint
	var season uint
	q := model.SearchQuery{
		Query:  params.Q,
		Type:   hint,
		Limit:  limit,
		Strong: strong,
	}

	if hint == media.Movies {
		if params.Year != nil {
			year = uint(*params.Year)
			q.Year = &year
		}

		if params.Season != nil {
			season = uint(*params.Season)
			q.Season = &season
		}
	}

	return q
}
func (s *Server) searchTorrents(params torrents.SearchTorrentsParams, key *models.Principal) middleware.Responder {
	l := s.log.WithField("req", "searchTorrents").WithField("key", key.Token).WithField("q", params.Q)
	l.Debug("Request")

	q := searchQueryFromParams(&params)
	mov, err := s.Torrents.Search(params.HTTPRequest.Context(), q)
	if err != nil {
		l.Errorf("Request failed: %s", err)
		return torrents.NewSearchTorrentsInternalServerError()
	}
	l.Debugf("Got %d results", len(mov))

	payload := torrents.SearchTorrentsOKBody{Results: []*models.SearchTorrentsResult{}}
	for i := range mov {
		payload.Results = append(payload.Results, convertTorrent(&mov[i], q.Type))
	}

	return torrents.NewSearchTorrentsOK().WithPayload(&payload)
}

func (s *Server) downloadTorrent(params torrents.DownloadTorrentParams, key *models.Principal) middleware.Responder {
	l := s.log.WithField("req", "downloadTorrent").WithField("key", key.Token).WithField("link", params.Link)
	l.Debug("Request")
	data, err := s.Torrents.Download(params.HTTPRequest.Context(), params.Link)
	if err != nil {
		l.Errorf("Request failed: %s", err)
		if errors.Is(err, torrents2.ErrExpiredDownloadLink) {
			return torrents.NewDownloadTorrentNotFound()
		}
		return torrents.NewDownloadTorrentInternalServerError()
	}
	l.Debug("Downloaded")
	rd := bytes.NewReader(data)
	return torrents.NewDownloadTorrentOK().WithPayload(io.NopCloser(rd))
}
