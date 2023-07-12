package pipeline

import (
	"context"
	"golang.org/x/time/rate"
)

type Handler func() (interface{}, error)

type Pipeline interface {
	Do(ctx context.Context, h Handler) (interface{}, error)
}

type Settings struct {
	Id    string
	Limit *rate.Limiter
}

func Open(settings Settings) Pipeline {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()

	p, ok := ctx.pipelines[settings.Id]
	if !ok {
		p = newPipeline(settings)
		ctx.pipelines[settings.Id] = p
	}

	return p
}

func Stop() {
	for _, p := range ctx.pipelines {
		p.stop()
	}
	ctx.pipelines = map[string]*pipeline{}
}
