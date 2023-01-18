package movies

import (
	"context"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/model"
)

func (s *service) Search(ctx context.Context, query string, limit uint) ([]model.Movie, error) {
	result, err := s.mainProvider.SearchMovies(ctx, query)
	if err != nil {
		return nil, err
	}

	if limit > 0 && uint(len(result)) > limit {
		result = result[:limit]
	}

	return result, nil
}
