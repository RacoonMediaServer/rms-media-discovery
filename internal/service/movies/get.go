package movies

import (
	"context"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/model"
)

func (s *Service) Get(ctx context.Context, id string) (*model.Movie, error) {
	return s.mainProvider.GetMovieInfo(ctx, id)
}
