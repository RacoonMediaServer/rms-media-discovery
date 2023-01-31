package thepiratebay

import (
	"context"
	"fmt"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/utils"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/pkg/model"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/pkg/navigator"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/pkg/provider"
	"github.com/PuerkitoBio/goquery"
	"github.com/apex/log"
	"net/url"
	"regexp"
	"strconv"
)

type tpbProvider struct {
	n navigator.Navigator
	l *log.Entry
}

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

func getFilter(hint model.SearchTypeHint) string {
	switch hint {
	case model.SearchType_Movies:
		return "video=on"
	case model.SearchType_Music:
		return "audio=on"
	case model.SearchType_Books:
		return "other=on"
	default:
		return "all=on"
	}
}

func (t *tpbProvider) ID() string {
	return "thepiratebay"
}

func (t *tpbProvider) SearchTorrents(ctx context.Context, q model.SearchQuery) ([]model.Torrent, error) {
	l := utils.LogFromContext(ctx, "thepiratebay", t.l)
	if t.n == nil {
		n, err := navigator.New()
		if err != nil {
			return []model.Torrent{}, fmt.Errorf("cannot create headless browser: %w", err)
		}
		t.n = n
	}
	p, err := t.n.NewPage(l, ctx)
	if err != nil {
		return []model.Torrent{}, fmt.Errorf("open page failed: %w", err)
	}
	defer p.Close()

	u := fmt.Sprintf("https://thepiratebay.org/search.php?q=%s&%s&search=Pirate+Search&page=0&orderby=",
		url.QueryEscape(q.Query),
		getFilter(q.Hint))

	err = p.Batch("searching").
		Goto(u).
		FetchContent().
		Error()

	if err != nil {
		return []model.Torrent{}, err
	}

	result := []model.Torrent{}

	doc := p.Document()
	doc.Find("#st").Each(func(i int, selection *goquery.Selection) {
		t := model.Torrent{}
		t.Title = selection.Find("span.list-item.item-name.item-title").First().Text()
		seeders, _ := strconv.ParseInt(selection.Find("span.list-item.item-seed").Text(), 10, 32)
		t.Seeders = uint(seeders)
		t.SizeMB = parseTorrentSize(selection.Find("span.list-item.item-size").Text())
		dl, ok := selection.Find("span.item-icons > a").Attr("href")
		if ok {
			t.Link = dl
		}

		if t.Link != "" {
			result = append(result, t)
		}
	})

	utils.SortTorrents(result)
	utils.Bound(result, q.Limit)

	for i := range result {
		torrent := &result[i]
		torrent.Downloader = func(ctx context.Context) ([]byte, error) {
			return []byte(torrent.Link), nil
		}
	}

	return result, nil
}

func New() provider.TorrentsProvider {
	return &tpbProvider{l: log.WithField("from", "thepiratebay")}
}
