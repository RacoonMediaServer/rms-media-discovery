package movies

import (
	"context"
	model2 "git.rms.local/RacoonMediaServer/rms-media-discovery/pkg/model"
	provider2 "git.rms.local/RacoonMediaServer/rms-media-discovery/pkg/provider"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/pkg/provider/imdb"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/pkg/provider/kinopoisk"
)

type Service interface {
	Search(ctx context.Context, query string, limit uint) ([]model2.Movie, error)
}

type service struct {
	mainProvider provider2.MovieInfoProvider
}

func New(access model2.AccessProvider) Service {
	return &service{
		mainProvider: provider2.NewFallbackProvider([]provider2.MovieInfoProvider{
			imdb.NewProvider(access),
			kinopoisk.NewKinopoiskProvider(access),
		}),
	}
}
