package movies

import (
	"context"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/model"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/provider"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/provider/imdb"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/provider/kinopoisk"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/provider/tmdb"
)

type Service interface {
	Search(ctx context.Context, query string, limit uint) ([]model.Movie, error)
	Get(ctx context.Context, id string) (*model.Movie, error)
}

type service struct {
	mainProvider provider.MovieInfoProvider
}

func New(access model.AccessProvider) Service {
	return &service{
		mainProvider: provider.NewFallbackProvider([]provider.MovieInfoProvider{
			tmdb.NewProvider(access),
			kinopoisk.NewProvider(access),
			imdb.NewProvider(access),
		}),
	}
}
