package server

import (
	"context"
	"github.com/RacoonMediaServer/rms-media-discovery/internal/server/models"
	"github.com/RacoonMediaServer/rms-media-discovery/internal/server/restapi/operations/music"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/model"
	"github.com/go-openapi/runtime/middleware"
)

func toPointer[T any](s T) *T {
	return &s
}

func convertSearchType(val *string) model.MusicSearchType {
	if val == nil {
		return model.SearchAny
	}
	switch *val {
	case "artist":
		return model.SearchArtist
	case "album":
		return model.SearchAlbum
	case "track":
		return model.SearchTrack
	default:
		return model.SearchAny
	}
}

func convertSearchMusicResult(m model.Music) *models.SearchMusicResult {
	result := models.SearchMusicResult{
		Title: toPointer(m.Title()),
	}
	if m.IsArtist() {
		a := m.AsArtist()
		result.Artist = a.Name
		result.Type = toPointer("artist")
		result.Picture = a.PictureUrl
		result.AlbumsCount = int64(a.Albums)
	} else if m.IsAlbum() {
		a := m.AsAlbum()
		result.Artist = a.Artist
		result.Album = a.Title
		result.Type = toPointer("album")
		result.Picture = a.CoverUrl
		result.Genres = a.Genres
		result.TracksCount = int64(a.Tracks)
		result.ReleaseYear = 0 // TODO
	} else if m.IsTrack() {
		t := m.AsTrack()
		result.Artist = t.Artist
		result.Album = t.Album
		result.Type = toPointer("track")
		result.Picture = t.CoverUrl
	}

	return &result
}

func (s *Server) searchMusic(params music.SearchMusicParams, key *models.Principal) middleware.Responder {
	l := s.log.WithField("key", key.Token).WithField("req", "searchMusic").WithField("q", params.Q)
	l.Debug("Request")

	var limit uint
	if params.Limit != nil {
		limit = uint(*params.Limit)
	}
	searchType := convertSearchType(params.Type)
	result, err := s.Music.Search(context.WithValue(params.HTTPRequest.Context(), "log", l), params.Q, limit, searchType)
	if err != nil {
		l.Errorf("Request failed: %s", err)
		return music.NewSearchMusicInternalServerError()
	}

	payload := music.SearchMusicOKBody{
		Results: make([]*models.SearchMusicResult, len(result)),
	}

	for i := range result {
		payload.Results[i] = convertSearchMusicResult(result[i])
	}

	l.Debugf("Got %d results", len(result))

	return music.NewSearchMusicOK().WithPayload(&payload)
}
