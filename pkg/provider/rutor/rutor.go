package rutor

import (
	"context"
	"errors"
	"fmt"
	"github.com/RacoonMediaServer/rms-media-discovery/internal/utils"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/media"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/model"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/provider"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/requester"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/scraper"
	"net/url"
)

const rutorDomain = "rutor.info"

type rutorProvider struct {
}

func (r rutorProvider) ID() string {
	return "rutor"
}

func composeURL(q model.SearchQuery) string {
	u := "http://" + rutorDomain + "/search/"
	if q.Type == media.Movies {
		if q.Year != nil || q.Season != nil {
			u += "0/0/100/2/" // более релевантная сортировка для наших задач
		}
		if q.Year != nil {
			q.Query += fmt.Sprintf(" %d", *q.Year)
		}
		if q.Season != nil {
			q.Query += fmt.Sprintf(" S%02d", *q.Season)
		}
	}

	u += url.PathEscape(q.Query)

	return u
}

func (r rutorProvider) SearchTorrents(ctx context.Context, q model.SearchQuery) ([]model.Torrent, error) {
	l := utils.LogFromContext(ctx, r.ID())
	c := scraper.New("rutor")
	c.SetContext(ctx)

	var result []model.Torrent
	available := false

	u := composeURL(q)
	err := c.
		Select(`#index > table > tbody > tr`, r.torrentsParser(&result)).
		Select(`#logo`, pageChecker(&available)). // проверяем, что видим реально rutor, а не заглушку РКН
		Get(u)
	if err != nil {
		return result, err
	}
	if !available {
		return nil, errors.New("domain is unavailable")
	}

	l.Debugf("Got %d results", len(result))

	return result, nil
}

func NewProvider() provider.TorrentsProvider {
	return &rutorProvider{}
}

func (r rutorProvider) newDownloadLink(url string) model.DownloadFunc {
	return func(ctx context.Context) ([]byte, error) {
		r := requester.New(r)
		data, contentType, err := r.Download(ctx, url)
		if contentType != "application/x-bittorrent" {
			return nil, errors.New("not accepted content-type received")
		}
		return data, err
	}
}
