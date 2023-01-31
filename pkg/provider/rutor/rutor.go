package rutor

import (
	"context"
	"errors"
	"fmt"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/utils"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/pkg/model"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/pkg/provider"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/pkg/requester"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/pkg/scraper"
	"net/url"
)

const domain = "rutor.info"

type rutorProvider struct {
}

func (r rutorProvider) ID() string {
	return "rutor"
}

func (r rutorProvider) SearchTorrents(ctx context.Context, q model.SearchQuery) ([]model.Torrent, error) {
	c := scraper.New("rutor")
	c.SetContext(ctx)

	var result []model.Torrent
	available := false

	u := fmt.Sprintf("http://%s/search/%s", domain, url.PathEscape(q.Query))
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

	utils.SortTorrents(result)
	result = utils.Bound(result, q.Limit)

	if q.Detailed {
		r.parseDetails(c, result)
	}
	return result, nil
}

func (r rutorProvider) parseDetails(c scraper.Scraper, torrents []model.Torrent) {
	c = c.Clone()
	sel := c.Select("#details > tbody > tr:nth-child(1) > td:nth-child(2)", detailsParser)
	for i := range torrents {
		t := &torrents[i]
		sel.GetAsync("http://"+domain+t.DetailLink, t)
	}

	c.Wait()
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
