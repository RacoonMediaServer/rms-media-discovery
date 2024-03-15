package music

import (
	"context"
	"github.com/RacoonMediaServer/rms-media-discovery/internal/utils"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/model"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/provider"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/provider/deezer"
)

type Service struct {
	mainProvider provider.MusicInfoProvider
}

const maxResults = 10

func (s Service) Search(ctx context.Context, query string, limit uint, searchType model.MusicSearchType) ([]model.Music, error) {
	if limit > maxResults || limit == 0 {
		limit = maxResults
	}

	result, err := s.mainProvider.SearchMusic(ctx, query, limit, searchType)

	result = utils.Bound(result, limit)
	return result, err
}

func New() *Service {
	return &Service{
		mainProvider: deezer.NewProvider(),
	}
}
