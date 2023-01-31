package thepiratebay

import (
	"context"
	"fmt"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/utils"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/pkg/model"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/pkg/navigator"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/pkg/provider"
	"github.com/apex/log"
	"net/url"
)

type tpbProvider struct {
	n navigator.Navigator
	l *log.Entry
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
	l := utils.LogFromContext(ctx, t.ID(), t.l)
	if t.n == nil {
		n, err := navigator.New(t.ID())
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
	p.Document().Find("#st").Each(torrentsParser(&result))

	utils.SortTorrents(result)
	result = utils.Bound(result, q.Limit)

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
