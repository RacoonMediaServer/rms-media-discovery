package rutor

import (
	"context"
	"fmt"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/model"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/provider"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/requester"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/utils"
	"github.com/apex/log"
	"github.com/gocolly/colly/v2"
	"net/url"
	"strconv"
	"sync"
)

const domain = "rutor.info"

type rutorProvider struct {
	log *log.Entry
}

func (r rutorProvider) ID() string {
	return "rutor"
}

func (r rutorProvider) SearchTorrents(ctx context.Context, query string, limit uint) ([]model.Torrent, error) {
	l := utils.LogFromContext(ctx, "rutor", r.log)
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36"),
		colly.AllowURLRevisit(),
	)
	utils.CollyWithContext(c, ctx)

	var result []model.Torrent

	c.OnHTML("#index > table > tbody > tr", func(e *colly.HTMLElement) {
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
	})
	u := fmt.Sprintf("http://%s/search/%s", domain, url.PathEscape(query))
	if err := c.Visit(u); err != nil {
		return nil, err
	}
	c.Wait()

	utils.SortTorrents(result)
	result = utils.Bound(result, limit)

	r.parseDetails(l, c, result)
	return result, nil
}

func (r rutorProvider) parseDetails(l *log.Entry, c *colly.Collector, torrents []model.Torrent) {
	wg := sync.WaitGroup{}
	for i := range torrents {
		t := &torrents[i]
		c := c.Clone()

		wg.Add(1)
		go func() {
			defer wg.Done()
			u := fmt.Sprintf("http://%s%s", domain, t.DetailLink)
			if err := c.Visit(u); err != nil {
				l.Warnf("Extract details failed: %s", err)
			}
			c.Wait()
		}()
	}
	wg.Wait()
}

func NewProvider() provider.TorrentsProvider {
	return &rutorProvider{
		log: log.WithField("from", "rutor"),
	}
}

func (r rutorProvider) newDownloadLink(url string) model.DownloadFunc {
	return func(ctx context.Context) ([]byte, error) {
		r := requester.New(r)
		return r.Download(ctx, url)
	}
}
