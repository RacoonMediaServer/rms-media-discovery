package movies

import (
	"context"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/model"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/utils"
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
