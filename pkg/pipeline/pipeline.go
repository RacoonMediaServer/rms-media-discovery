package pipeline

import (
	"context"
	"golang.org/x/time/rate"
	"sync"
)

const MaxAttempts = 20

type pipeline struct {
	settings Settings
	ch       chan *task
	wg       sync.WaitGroup
}

type callResult struct {
	data interface{}
	err  error
}

var ctx struct {
	pipelines map[string]*pipeline
	mu        sync.Mutex
}

func init() {
	ctx.pipelines = make(map[string]*pipeline)
}

func newPipeline(settings Settings) *pipeline {
	p := &pipeline{
		settings: settings,
		ch:       make(chan *task),
	}
	if p.settings.Limit == nil {
		p.settings.Limit = rate.NewLimiter(rate.Inf, 1)
	}

	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		p.process()
	}()

	return p
}

type task struct {
	h  Handler
	ch chan callResult
}

func (p *pipeline) Do(ctx context.Context, h Handler) (interface{}, error) {
	t := &task{
		h:  h,
		ch: make(chan callResult, 2),
	}

	select {
	case p.ch <- t:
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case r := <-t.ch:
		return r.data, r.err
	}
}

func (p *pipeline) process() {
	for {
		t, ok := <-p.ch
		if !ok {
			return
		}

		_ = p.settings.Limit.Wait(context.Background())
		data, err := t.h()
		t.ch <- callResult{data: data, err: err}
		close(t.ch)
	}
}

func (p *pipeline) stop() {
	close(p.ch)
	p.wg.Wait()
}
