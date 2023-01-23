package movies

import (
	"context"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/model"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/provider"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/provider/imdb"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/provider/kinopoisk"
)

type Service interface {
	Search(ctx context.Context, query string, limit uint) ([]model.Movie, error)
}

type service struct {
	mainProvider provider.MovieInfoProvider
}

func New(access model.AccessProvider) Service {
	return &service{
		mainProvider: provider.NewFallbackProvider([]provider.MovieInfoProvider{
			imdb.NewProvider(access),
			kinopoisk.NewKinopoiskProvider(access),
		}),
	}
}
