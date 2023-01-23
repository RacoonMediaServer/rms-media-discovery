package rutracker

import (
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/model"
	"github.com/gocolly/colly/v2"
	"regexp"
	"strconv"
)

var (
	captchaSidExpr  = regexp.MustCompile(`<input[^>]*name="cap_sid"[^>]*value="([^"]+)"[^>]*>`)
	captchaCodeExpr = regexp.MustCompile(`<input[^>]*name="(cap_code_[^"]+)"[^>]*value="[^"]*"[^>]*>`)
	captchaUrlExpr  = regexp.MustCompile(`<img[^>]*src="([^"]+\/captcha\/[^"]+)"[^>]*>`)
	extractSizeExpr = regexp.MustCompile(`^(\d+(.\d+)?) (MB|GB)`)
)

func parseTorrentSize(text string) float32 {
	matches := extractSizeExpr.FindStringSubmatch(text)
	if matches != nil {
		result, err := strconv.ParseFloat(matches[1], 32)
		if err != nil {
			return 0
		}
		if matches[3] == "GB" {
			result *= 1024.
		}
		return float32(result)
	}

	return 0
}

func parseTorrent(e *colly.HTMLElement) model.Torrent {
	torrent := model.Torrent{}
	torrent.Title = e.DOM.Find(`a.tLink`).Text()

	dl := e.DOM.Find(`a.tr-dl`)
	link, _ := dl.Attr("href")
	torrent.Link = link
	torrent.SizeMB = parseTorrentSize(dl.Text())

	seeds := e.DOM.Find(`b.seedmed`).Text()
	seedersCount, _ := strconv.ParseUint(seeds, 10, 32)
	torrent.Seeders = uint(seedersCount)

	leechs := e.DOM.Find(`td.leechmed`).Text()
	peers, _ := strconv.Atoi(leechs)
	torrent.Seeders += uint(peers)

	return torrent
}
