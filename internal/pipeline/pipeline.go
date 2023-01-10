package pipeline

import (
	"context"
	"sync"
)

const maxAttempts = 20

type pipeline struct {
	settings Settings
	ch       chan *task
	wg       sync.WaitGroup
}

var ctx struct {
	pipelines map[string]*pipeline
	mu        sync.Mutex
}

func newPipeline(settings Settings) *pipeline {
	p := &pipeline{
		settings: settings,
		ch:       make(chan *task),
	}

	if settings.MaxAttempts == 0 {
		settings.MaxAttempts = maxAttempts
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
	ch chan Result
}

func (p *pipeline) Do(ctx context.Context, h Handler) (interface{}, error) {
	t := &task{
		h:  h,
		ch: make(chan Result, 2),
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
		return r.Result, r.Err
	}
}

func (p *pipeline) process() {
	for {
		t, ok := <-p.ch
		if !ok {
			return
		}

		r := Result{}
		for a := 0; a < int(p.settings.MaxAttempts); a++ {
			r = t.h()
			if r.Done {
				break
			}
		}
		t.ch <- r
		close(t.ch)
	}
}

func (p *pipeline) stop() {
	close(p.ch)
	p.wg.Wait()
}
