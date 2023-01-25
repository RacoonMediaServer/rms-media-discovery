package torrents

import (
	"context"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/model"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/provider"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/utils"
	"sync"
)

type aggregator struct {
	providers []provider.TorrentsProvider
}

func (a aggregator) SearchTorrents(ctx context.Context, query string, limit uint) ([]model.Torrent, error) {
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
			r, err := p.SearchTorrents(ctx, query, limit)
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

	utils.SortTorrents(total)
	return utils.Bound(total, limit), nil
}

func newAggregator(providers []provider.TorrentsProvider) provider.TorrentsProvider {
	return &aggregator{providers: providers}
}
