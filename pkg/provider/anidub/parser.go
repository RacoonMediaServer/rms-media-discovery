package anidub

import (
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/model"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/scraper"
	"github.com/gocolly/colly/v2"
	"regexp"
	"strconv"
)

var extractSizeExpr = regexp.MustCompile(`(\d+(.\d+)?).(MB|GB)`)

func parseTorrentSize(text string) uint64 {
	matches := extractSizeExpr.FindStringSubmatch(text)
	if matches != nil {
		result, err := strconv.ParseFloat(matches[1], 32)
		if err != nil {
			return 0
		}
		if matches[3] == "GB" {
			result *= 1024.
		}
		return uint64(result)
	}

	return 0
}

func (a *anidubProvider) searchItemsParser(result *[]string) scraper.HTMLCallback {
	return func(e *colly.HTMLElement, userData interface{}) {
		link := e.ChildAttr("h3:nth-child(1) > a:nth-child(1)", "href")
		a.l.Debugf("Found item: %s", link)
		*result = append(*result, link)
	}
}

func titleParser(t *model.Torrent) scraper.HTMLCallback {
	return func(e *colly.HTMLElement, userData interface{}) {
		t.Title = e.Text
	}
}

func linkParser(t *model.Torrent) scraper.HTMLCallback {
	return func(e *colly.HTMLElement, userData interface{}) {
		t.Link = e.Attr("href")
	}
}

func metricsParser(t *model.Torrent) scraper.HTMLCallback {
	return func(e *colly.HTMLElement, userData interface{}) {
		t.SizeMB = parseTorrentSize(e.ChildText(`span.red`))
		seeders, err := strconv.ParseUint(e.ChildText(`.li_download_m`), 10, 32)
		if err == nil {
			t.Seeders = uint(seeders)
		}
	}
}
