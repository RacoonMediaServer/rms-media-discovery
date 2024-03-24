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
	"net/url"
)

type anidubProvider struct {
	l *log.Entry
}

func (a *anidubProvider) ID() string {
	return "anidub"
}

func composeSearchParams(query model.SearchQuery) url.Values {
	return url.Values{
		"do":           []string{"search"},
		"subaction":    []string{"search"},
		"search_start": []string{"1"},
		"search_full":  []string{"1"},
		"result_from":  []string{"1"},
		"story":        []string{query.Query},
		"titleonly":    []string{"3"},
		"searchuser":   []string{""},
		"replyless":    []string{"0"},
		"replylimit":   []string{"0"},
		"searchdate":   []string{"0"},
		"beforeafter":  []string{"after"},
		"sortby":       []string{""},
		"resorder":     []string{"desc"},
		"showposts":    []string{"1"},
		"catlist[]":    []string{"2", "13", "3", "4", "9"},
	}
}

func isMustSkip(q model.SearchQuery, info *animeInfo) bool {
	if q.Year != nil && *q.Year != info.Year {
		return true
	}
	if q.Season != nil && *q.Season != 1 {
		return true
	}
	return false
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
		t := model.Torrent{}
		info := animeInfo{}
		err = s.
			Select(`#news-title`, titleParser(&t)).
			Select(`.torrent_h > a:nth-child(1)`, linkParser(&t)).
			Select(`div.list:nth-child(2)`, metricsParser(&t)).
			Select(`.xfinfodata`, infoParser(&info)).
			Get(link)
		if err != nil {
			a.l.Errorf("Extract torrent info failed (%s): %s", link, err)
			continue
		}
		if t.IsValid() && !isMustSkip(query, &info) {
			t.Downloader = a.newDownloadLink("https://tr.anidub.com" + t.Link)
			results = append(results, t)
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
