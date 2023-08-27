package music

import (
	"context"
	"github.com/RacoonMediaServer/rms-media-discovery/internal/utils"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/model"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/provider"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/provider/deezer"
)

type service struct {
	mainProvider provider.MusicInfoProvider
}

const maxResults = 10

func (s service) Search(ctx context.Context, query string, limit uint) ([]model.Music, error) {
	if limit > maxResults || limit == 0 {
		limit = maxResults
	}

	result, err := s.mainProvider.SearchMusic(ctx, query, limit)

	result = utils.Bound(result, limit)
	return result, err
}

type Service interface {
	Search(ctx context.Context, query string, limit uint) ([]model.Music, error)
}

func New() Service {
	return &service{
		mainProvider: deezer.NewProvider(),
	}
}