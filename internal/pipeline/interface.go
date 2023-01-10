package pipeline

import "context"

type Result struct {
	Done   bool
	Err    error
	Result interface{}
}

type Handler func() Result

type Pipeline interface {
	Do(ctx context.Context, h Handler) (interface{}, error)
}

type Settings struct {
	Id          string
	MaxAttempts uint
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
}
