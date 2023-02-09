package rutor

import (
	"git.rms.local/RacoonMediaServer/rms-media-discovery/pkg/model"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/pkg/scraper"
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

func (r rutorProvider) torrentsParser(result *[]model.Torrent) scraper.HTMLCallback {
	return func(e *colly.HTMLElement, userData interface{}) {
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
			*result = append(*result, t)
		}
	}
}

func pageChecker(isOk *bool) scraper.HTMLCallback {
	return func(e *colly.HTMLElement, userData interface{}) {
		*isOk = true
	}
}

func detailsParser(e *colly.HTMLElement, userData interface{}) {
	// TODO: либо убрать, либо парсить содержимое
	//t := userData.(*model.Torrent)
	//parser := heuristic.MediaInfoParser{}
	//t.Media = parser.Parse(e.Text)
}
