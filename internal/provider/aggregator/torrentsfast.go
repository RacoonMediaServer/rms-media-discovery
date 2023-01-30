package aggregator

import (
	"context"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/model"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/provider"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/utils"
	"sync"
)

type torrentsFastAggregator struct {
	providers []provider.TorrentsProvider
}

func (a torrentsFastAggregator) ID() string {
	return "aggregator"
}

func (a torrentsFastAggregator) SearchTorrents(ctx context.Context, q model.SearchQuery) ([]model.Torrent, error) {
	type result struct {
		torrents []model.Torrent
		err      error
	}

	ch := make(chan result, len(a.providers))
	defer close(ch)

	wg := sync.WaitGroup{}
	defer wg.Wait()

	waitCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	wg.Add(len(a.providers))
	for _, p := range a.providers {
		go func(p provider.TorrentsProvider) {
			defer wg.Done()
			r, err := p.SearchTorrents(waitCtx, q)
			ch <- result{torrents: r, err: err}
		}(p)
	}

	var total []model.Torrent
	anySuccess := false
	var lastErr error
	for _ = range a.providers {
		r := <-ch
		if r.err == nil {
			total = append(total, r.torrents...)
			anySuccess = true
			if uint(len(total)) >= q.Limit {
				break
			}
		} else {
			lastErr = r.err
		}
	}

	if !anySuccess {
		return []model.Torrent{}, lastErr
	}

	utils.SortTorrents(total)
	return utils.Bound(total, q.Limit), nil
}
