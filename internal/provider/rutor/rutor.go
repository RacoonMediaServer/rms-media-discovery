package rutor

import (
	"context"
	"fmt"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/model"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/provider"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/requester"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/scraper"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/utils"
	"github.com/apex/log"
	"github.com/gocolly/colly/v2"
	"net/url"
	"strconv"
)

const domain = "rutor.info"

type rutorProvider struct {
	log *log.Entry
}

func (r rutorProvider) ID() string {
	return "rutor"
}

func (r rutorProvider) SearchTorrents(ctx context.Context, query string, limit uint) ([]model.Torrent, error) {
	c := scraper.New("rutor")
	c.SetContext(ctx)

	var result []model.Torrent

	u := fmt.Sprintf("http://%s/search/%s", domain, url.PathEscape(query))
	err := c.Select("#index > table > tbody > tr", func(e *colly.HTMLElement, userData interface{}) {
		downloadLink := e.ChildAttr("td:nth-child(2) > a.downgif", "href")
		title := e.ChildText("td:nth-child(2) > a:nth-child(3)")
		scrapLink := e.ChildAttr("td:nth-child(2) > a:nth-child(3)", "href")
		size := parseTorrentSize(e.Text)
		seeds, _ := strconv.ParseUint(e.ChildText("td > span.green"), 10, 32)

		if downloadLink != "" {
			t := model.Torrent{
				Title:      title,
				SizeMB:     size,
				Seeders:    uint(seeds),
				DetailLink: scrapLink,
				Downloader: r.newDownloadLink(downloadLink),
			}
			result = append(result, t)
		}
	}).Get(u)
	if err != nil {
		return nil, err
	}

	utils.SortTorrents(result)
	result = utils.Bound(result, limit)

	r.parseDetails(c, result)
	return result, nil
}

func (r rutorProvider) parseDetails(c scraper.Scraper, torrents []model.Torrent) {
	c = c.Clone()
	sel := c.Select("#logo > a > img", func(e *colly.HTMLElement, userData interface{}) {
		// TODO: scrap
	})
	for i := range torrents {
		t := &torrents[i]
		sel.GetAsync("http://"+domain+t.DetailLink, t)
	}

	c.Wait()
}

func NewProvider() provider.TorrentsProvider {
	return &rutorProvider{
		log: log.WithField("from", "rutor"),
	}
}

func (r rutorProvider) newDownloadLink(url string) model.DownloadFunc {
	return func(ctx context.Context) ([]byte, error) {
		r := requester.New(r)
		data, _, err := r.Download(ctx, url)
		return data, err
	}
}
