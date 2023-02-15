package movies

import (
	"context"
	"github.com/RacoonMediaServer/rms-media-discovery/internal/utils"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/model"
)

const maxResults uint = 10

func (s *service) Search(ctx context.Context, query string, limit uint) ([]model.Movie, error) {
	if limit == 0 || limit > maxResults {
		limit = maxResults
	}

	result, err := s.mainProvider.SearchMovies(ctx, query, limit)
	if err != nil {
		return nil, err
	}

	result = utils.Bound(result, limit)

	return result, nil
}
