package thepiratebay

import (
	"git.rms.local/RacoonMediaServer/rms-media-discovery/pkg/model"
	"github.com/PuerkitoBio/goquery"
	"regexp"
	"strconv"
)

var (
	extractSizeExpr = regexp.MustCompile(`^(^\d+(\.\d+)?).(G|M)iB`)
)

func parseTorrentSize(text string) uint64 {
	matches := extractSizeExpr.FindStringSubmatch(text)
	if matches != nil {
		result, err := strconv.ParseFloat(matches[1], 32)
		if err != nil {
			return 0
		}
		if matches[3] == "G" {
			result *= 1024.
		}
		return uint64(result)
	}

	return 0
}

func torrentsParser(results *[]model.Torrent) func(int, *goquery.Selection) {
	return func(i int, selection *goquery.Selection) {
		t := model.Torrent{}
		t.Title = selection.Find("span.list-item.item-name.item-title").First().Text()
		seeders, _ := strconv.ParseInt(selection.Find("span.list-item.item-seed").Text(), 10, 32)
		t.Seeders = uint(seeders)
		t.SizeMB = parseTorrentSize(selection.Find("span.list-item.item-size").Text())
		dl, ok := selection.Find("span.item-icons > a").Attr("href")
		if ok {
			t.Link = dl
		}

		if t.Link != "" && t.Title != "No results returned" {
			*results = append(*results, t)
		}
	}
}
