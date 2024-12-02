package rutracker

import (
	"regexp"
	"strconv"

	"github.com/RacoonMediaServer/rms-media-discovery/pkg/model"
	"github.com/gocolly/colly/v2"
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

func parseTorrent(e *colly.HTMLElement) model.Torrent {
	torrent := model.Torrent{}
	torrent.Title = e.ChildText(`a.tLink`)

	torrent.Link = e.ChildAttr(`a.tr-dl`, "href")
	torrent.SizeMB = parseTorrentSize(e.ChildText(`a.tr-dl`))

	seeds := e.ChildText(`b.seedmed`)
	seedersCount, _ := strconv.ParseUint(seeds, 10, 32)
	torrent.Seeders = uint(seedersCount)

	leechs := e.ChildText(`td.leechmed`)
	peers, _ := strconv.Atoi(leechs)
	torrent.Seeders += uint(peers)

	torrent.DetailLink = e.ChildAttr(`a.tr-dl`, "href")

	return torrent
}
