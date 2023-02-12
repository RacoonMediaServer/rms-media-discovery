package thepiratebay

import (
	"context"
	"fmt"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/utils"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/pkg/media"
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

func getFilter(hint media.ContentType) string {
	switch hint {
	case media.Movies:
		return "video=on"
	case media.Music:
		return "audio=on"
	case media.Books:
		return "other=on"
	default:
		return "all=on"
	}
}

func (t *tpbProvider) ID() string {
	return "thepiratebay"
}
func applySearchHints(q *model.SearchQuery) {
	// применяем дополнительные параметры поиска так, как это лучше всего будет работать на конкретном трекере
	if q.Type == media.Movies {
		if q.Year != nil {
			q.Query += fmt.Sprintf(" %d", *q.Year)
		}
		if q.Season != nil {
			q.Query += fmt.Sprintf(" S%02d", *q.Season)
		}
		q.Query += " rus" // хотим русскую озвучку
	}
}
func (t *tpbProvider) SearchTorrents(ctx context.Context, q model.SearchQuery) ([]model.Torrent, error) {
	applySearchHints(&q)
	l := utils.LogFromContext(ctx, t.ID())
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
		getFilter(q.Type))

	err = p.Batch("searching").
		Goto(u).
		FetchContent().
		Error()

	if err != nil {
		return []model.Torrent{}, err
	}

	result := []model.Torrent{}
	p.Document().Find("#st").Each(torrentsParser(&result))

	for i := range result {
		torrent := &result[i]
		torrent.Downloader = func(ctx context.Context) ([]byte, error) {
			return []byte(torrent.Link), nil
		}
	}
	l.Debugf("Got %d results", len(result))

	return result, nil
}

func New() provider.TorrentsProvider {
	return &tpbProvider{}
}
