package anidub

import (
	"context"
	"errors"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/media"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/model"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/provider"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/requester"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/scraper"
	"github.com/apex/log"
)

type anidubProvider struct {
	l *log.Entry
}

func (a *anidubProvider) ID() string {
	return "anidub"
}

func composeSearchParams(query model.SearchQuery) map[string]string {
	return map[string]string{
		"do":           "search",
		"subaction":    "search",
		"search_start": "1",
		"search_full":  "1",
		"result_from":  "1",
		"story":        query.Query,
		"titleonly":    "3",
		"searchuser":   "",
		"replyless":    "0",
		"replylimit":   "0",
		"searchdate":   "0",
		"beforeafter":  "after",
		"sortby":       "",
		"resorder":     "desc",
		"showposts":    "1",
		"catlist[]":    "0",
	}
}

func (a *anidubProvider) SearchTorrents(ctx context.Context, query model.SearchQuery) ([]model.Torrent, error) {
	if query.Type != media.Movies && query.Type != media.Other {
		return []model.Torrent{}, nil
	}
	s := scraper.New(a.ID())
	s.SetContext(ctx)

	var links []string
	err := s.Select(`div[class="dpad searchitem"]`, a.searchItemsParser(&links)).Post("https://tr.anidub.com/index.php?do=search", composeSearchParams(query))
	if err != nil {
		return []model.Torrent{}, err
	}

	var results []model.Torrent
	for _, link := range links {
		t := &model.Torrent{}
		err = s.
			Select(`#news-title`, titleParser(t)).
			Select(`.torrent_h > a:nth-child(1)`, linkParser(t)).
			Select(`div.list:nth-child(2)`, metricsParser(t)).
			Get(link)
		if err != nil {
			a.l.Errorf("Extract torrent info failed (%s): %s", link, err)
			continue
		}
		if t.IsValid() {
			t.Downloader = a.newDownloadLink("https://tr.anidub.com" + t.Link)
			results = append(results, *t)
		}
	}
	return results, err
}

func New() provider.TorrentsProvider {
	return &anidubProvider{
		l: log.WithField("from", "anidub"),
	}
}

func (a *anidubProvider) newDownloadLink(url string) model.DownloadFunc {
	return func(ctx context.Context) ([]byte, error) {
		r := requester.New(a)
		data, contentType, err := r.Download(ctx, url)
		if contentType != "application/x-bittorrent" {
			return nil, errors.New("not accepted content-type received")
		}
		return data, err
	}
}
