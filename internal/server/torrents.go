package server

import (
	"errors"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/server/models"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/server/restapi/operations/torrents"
	torrents2 "git.rms.local/RacoonMediaServer/rms-media-discovery/internal/service/torrents"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/pkg/model"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"net/http"
)

func convertTorrent(t *model.Torrent) *models.SearchTorrentsResult {
	result := &models.SearchTorrentsResult{
		Link:    &t.Link,
		Seeders: int64(t.Seeders),
		Size:    int64(t.SizeMB),
		Title:   t.Title,
	}
	if t.Media != nil {
		result.Media = new(models.SearchTorrentsResultMedia)
		for _, v := range t.Media.Video {
			target := &models.SearchTorrentsResultMediaVideoItems0{
				AspectRatio: v.AspectRatio,
				Codec:       v.Codec,
				Height:      int64(v.Height),
				Width:       int64(v.Width),
			}
			result.Media.Video = append(result.Media.Video, target)
		}
		for _, a := range t.Media.Audio {
			target := &models.SearchTorrentsResultMediaAudioItems0{
				Codec:    a.Codec,
				Language: a.Language,
				Voice:    a.Voice,
			}
			result.Media.Audio = append(result.Media.Audio, target)
		}
		for _, s := range t.Media.Subtitle {
			target := &models.SearchTorrentsResultMediaSubtitlesItems0{
				Codec:    s.Codec,
				Language: s.Language,
			}
			result.Media.Subtitles = append(result.Media.Subtitles, target)
		}
	}

	return result
}

func searchQueryFromParams(params *torrents.SearchTorrentsParams) model.SearchQuery {
	var limit uint
	var detailed bool
	if params.Limit != nil {
		limit = uint(*params.Limit)
	}
	if params.Detailed != nil {
		detailed = *params.Detailed
	}

	hint := model.SearchType_Other
	if params.Type != nil {
		switch *params.Type {
		case "movies":
			hint = model.SearchType_Movies
		case "music":
			hint = model.SearchType_Music
		case "books":
			hint = model.SearchType_Books
		}
	}
	return model.SearchQuery{
		Query:    params.Q,
		Hint:     hint,
		Limit:    limit,
		Detailed: detailed,
	}
}
func (s *Server) searchTorrents(params torrents.SearchTorrentsParams, key *models.Principal) middleware.Responder {
	l := s.log.WithField("req", "searchTorrents").WithField("key", key.Token).WithField("q", params.Q)
	l.Debug("Request")

	mov, err := s.Torrents.Search(params.HTTPRequest.Context(), searchQueryFromParams(&params))
	if err != nil {
		l.Errorf("Request failed: %s", err)
		return torrents.NewSearchTorrentsInternalServerError()
	}
	l.Debugf("Got %d results", len(mov))

	payload := torrents.SearchTorrentsOKBody{Results: []*models.SearchTorrentsResult{}}
	for i := range mov {
		payload.Results = append(payload.Results, convertTorrent(&mov[i]))
	}

	return torrents.NewSearchTorrentsOK().WithPayload(&payload)
}

type downloadResponse struct {
	data []byte
}

func (r *downloadResponse) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {
	_, _ = rw.Write(r.data)
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
	return &downloadResponse{data: data}
}
