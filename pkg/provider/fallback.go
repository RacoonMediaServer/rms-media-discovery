package provider

import (
	"context"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/model"
)

type fallbackProvider struct {
	providers []MovieInfoProvider
}

func (f fallbackProvider) GetMovieInfo(ctx context.Context, id string) (result *model.Movie, err error) {
	for _, p := range f.providers {
		result, err = p.GetMovieInfo(ctx, id)
		if err == nil && result.Title != "" {
			return
		}
	}

	return
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
