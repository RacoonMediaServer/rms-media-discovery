package torrents

import (
	"context"
	"fmt"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/model"
	"sync"
)

type searchFunc func(ctx context.Context, q model.SearchQuery) ([]model.Torrent, error)

type searchTask struct {
	q model.SearchQuery
	f searchFunc

	ctx    context.Context
	cancel context.CancelFunc
	mu     sync.Mutex
	state  TaskStatus
}

func (t *searchTask) run() {
	defer func() {
		if rec := recover(); rec != nil {
			err, ok := rec.(error)
			if !ok {
				err = fmt.Errorf("%v", rec)
			}

			t.mu.Lock()
			defer t.mu.Unlock()
			t.state.Err = err
			t.state.Status = model.Error
		}
	}()

	results, err := t.f(t.ctx, t.q)

	t.mu.Lock()
	defer t.mu.Unlock()
	t.state.Results = results
	t.state.Err = err
	t.state.Status = model.Ready
	if err != nil {
		t.state.Status = model.Error
	}
}

func (t *searchTask) status() TaskStatus {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.state
}
