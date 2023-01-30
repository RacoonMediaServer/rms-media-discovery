package movies

import (
	"context"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/pkg/model"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/pkg/provider"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/pkg/provider/imdb"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/pkg/provider/kinopoisk"
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
