package torrents

import (
	"context"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/model"
)

func (s *service) Search(ctx context.Context, query string, hint SearchTypeHint, limit uint) ([]model.Torrent, error) {
	found, err := s.provider.SearchTorrents(ctx, query)
	if err != nil {
		return nil, err
	}
	if uint(len(found)) > limit {
		found = found[:limit]
	}
	return found, nil
}
