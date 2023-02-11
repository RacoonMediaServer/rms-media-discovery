package aggregator

import (
	"context"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/pkg/model"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/pkg/provider"
	"sync"
)

type torrentsPriorityAggregator struct {
	providers []provider.TorrentsProvider
}

func (a torrentsPriorityAggregator) ID() string {
	return "aggregator"
}

func (a torrentsPriorityAggregator) SearchTorrents(ctx context.Context, q model.SearchQuery) ([]model.Torrent, error) {
	type result struct {
		torrents []model.Torrent
		err      error
	}

	ch := make(chan result, len(a.providers))
	defer close(ch)

	wg := sync.WaitGroup{}
	wg.Add(len(a.providers))
	for _, p := range a.providers {
		go func(p provider.TorrentsProvider) {
			defer wg.Done()
			r, err := p.SearchTorrents(ctx, q)
			ch <- result{torrents: r, err: err}
		}(p)
	}
	wg.Wait()

	var total []model.Torrent
	anySuccess := false
	var lastErr error
	for _ = range a.providers {
		r := <-ch
		if r.err == nil {
			total = append(total, r.torrents...)
			anySuccess = true
		} else {
			lastErr = r.err
		}
	}

	if !anySuccess {
		return []model.Torrent{}, lastErr
	}

	return total, nil
}
