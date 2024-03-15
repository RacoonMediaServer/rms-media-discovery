package movies

import (
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/model"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/provider"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/provider/imdb"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/provider/kinopoisk"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/provider/tmdb"
)

type Service struct {
	mainProvider provider.MovieInfoProvider
}

func New(access model.AccessProvider) *Service {
	return &Service{
		mainProvider: provider.NewFallbackProvider([]provider.MovieInfoProvider{
			tmdb.NewProvider(access),
			kinopoisk.NewProvider(access),
			imdb.NewProvider(access),
		}),
	}
}
