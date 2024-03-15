package server

import (
	"bytes"
	"errors"
	"github.com/RacoonMediaServer/rms-media-discovery/internal/server/models"
	"github.com/RacoonMediaServer/rms-media-discovery/internal/server/restapi/operations/torrents"
	torrents2 "github.com/RacoonMediaServer/rms-media-discovery/internal/service/torrents"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/media"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/model"
	"github.com/go-openapi/runtime/middleware"
	"io"
)

const pollIntervalMs = 1000

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

func taskStatusToString(status model.TaskStatus) string {
	strings := map[model.TaskStatus]string{
		model.Working: "working",
		model.Ready:   "ready",
		model.Error:   "error",
	}
	return strings[status]
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

	switch hint {
	case media.Movies:
		if params.Year != nil {
			year = uint(*params.Year)
			q.Year = &year
		}

		if params.Season != nil {
			season = uint(*params.Season)
			q.Season = &season
		}
	case media.Music:
		if params.Discography != nil {
			q.Discography = *params.Discography
		}
	default:
	}

	return q
}

func searchQueryFromAsyncParams(params *torrents.SearchTorrentsAsyncBody) model.SearchQuery {
	var limit uint
	if params.Limit != 0 {
		limit = uint(params.Limit)
	}
	var strong bool
	if params.Strong != nil {
		strong = *params.Strong
	}

	hint := media.Other

	switch params.Type {
	case "movies":
		hint = media.Movies
	case "music":
		hint = media.Music
	case "books":
		hint = media.Books
	}

	var year uint
	var season uint
	q := model.SearchQuery{
		Query:  *params.Q,
		Type:   hint,
		Limit:  limit,
		Strong: strong,
	}

	switch hint {
	case media.Movies:
		if params.Year != 0 {
			year = uint(params.Year)
			q.Year = &year
		}

		if params.Season != 0 {
			season = uint(params.Season)
			q.Season = &season
		}
	case media.Music:
		if params.Discography != nil {
			q.Discography = *params.Discography
		}
	default:
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

func (s *Server) searchTorrentsAsync(params torrents.SearchTorrentsAsyncParams, key *models.Principal) middleware.Responder {
	l := s.log.WithField("req", "searchTorrentsAsync").WithField("key", key.Token).WithField("q", params.SearchParameters.Q)
	l.Debug("Request")

	q := searchQueryFromAsyncParams(&params.SearchParameters)
	taskID, err := s.Torrents.SearchAsync(q)
	if err != nil {
		l.Errorf("Run async task failed: %s", err)
		return torrents.NewSearchTorrentsAsyncInternalServerError()
	}

	payload := torrents.SearchTorrentsAsyncOKBody{
		ID:             taskID,
		PollIntervalMs: pollIntervalMs,
	}
	return torrents.NewSearchTorrentsAsyncOK().WithPayload(&payload)
}

func (s *Server) searchTorrentsAsyncStatus(params torrents.SearchTorrentsAsyncStatusParams, key *models.Principal) middleware.Responder {
	l := s.log.WithField("req", "searchTorrentsAsyncStatus").WithField("key", key.Token).WithField("id", params.ID)
	l.Debug("Request")

	status, err := s.Torrents.Status(params.ID)
	if err != nil {
		l.Errorf("Get async task failed: %s", err)
		if errors.Is(err, torrents2.ErrTaskNotFound) {
			return torrents.NewSearchTorrentsAsyncStatusNotFound()
		}
		return torrents.NewSearchTorrentsAsyncStatusInternalServerError()
	}

	statusString := taskStatusToString(status.Status)
	payload := torrents.SearchTorrentsAsyncStatusOKBody{
		Results: make([]*models.SearchTorrentsResult, 0, len(status.Results)),
		Status:  &statusString,
	}
	if status.Err == nil && status.Status == model.Ready {
		for i := range status.Results {
			payload.Results = append(payload.Results, convertTorrent(&status.Results[i], status.ContentType))
		}
	} else {
		payload.Error = status.Err.Error()
	}

	return torrents.NewSearchTorrentsAsyncStatusOK().WithPayload(&payload)
}

func (s *Server) searchTorrentsAsyncCancel(params torrents.SearchTorrentsAsyncCancelParams, key *models.Principal) middleware.Responder {
	l := s.log.WithField("req", "searchTorrentsAsyncCancel").WithField("key", key.Token).WithField("id", params.ID)
	l.Debug("Request")

	err := s.Torrents.Cancel(params.ID)
	if err != nil {
		l.Errorf("Cancel async task failed: %s", err)
		if errors.Is(err, torrents2.ErrTaskNotFound) {
			return torrents.NewSearchTorrentsAsyncCancelNotFound()
		}
		return torrents.NewSearchTorrentsAsyncCancelInternalServerError()
	}
	return torrents.NewSearchTorrentsAsyncCancelOK()
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
