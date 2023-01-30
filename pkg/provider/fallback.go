package provider

import (
	"context"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/pkg/model"
)

type fallbackProvider struct {
	providers []MovieInfoProvider
}

func (f fallbackProvider) ID() string {
	if len(f.providers) == 0 {
		return ""
	}

	return f.providers[0].ID()
}

func (f fallbackProvider) SearchMovies(ctx context.Context, query string, limit uint) (result []model.Movie, err error) {
	for _, p := range f.providers {
		result, err = p.SearchMovies(ctx, query, limit)
		if err == nil && len(result) != 0 {
			return
		}
	}

	return
}

func NewFallbackProvider(providers []MovieInfoProvider) MovieInfoProvider {
	return &fallbackProvider{providers: providers}
}
