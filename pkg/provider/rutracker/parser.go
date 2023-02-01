package rutracker

import (
	"git.rms.local/RacoonMediaServer/rms-media-discovery/pkg/model"
	"github.com/PuerkitoBio/goquery"
	"regexp"
	"strconv"
)

var (
	extractSizeExpr = regexp.MustCompile(`^(\d+(.\d+)?).(M|G)B`)
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

func parseTorrent(e *goquery.Selection) model.Torrent {
	torrent := model.Torrent{}
	torrent.Title = e.Find(`a.tLink`).Text()

	dl := e.Find(`a.tr-dl`)
	link, _ := dl.Attr("href")
	torrent.Link = link
	torrent.SizeMB = parseTorrentSize(dl.Text())

	seeds := e.Find(`b.seedmed`).Text()
	seedersCount, _ := strconv.ParseUint(seeds, 10, 32)
	torrent.Seeders = uint(seedersCount)

	leechs := e.Find(`td.leechmed`).Text()
	peers, _ := strconv.Atoi(leechs)
	torrent.Seeders += uint(peers)

	torrent.DetailLink, _ = e.Find(`a.tLink`).Attr("href")

	return torrent
}
