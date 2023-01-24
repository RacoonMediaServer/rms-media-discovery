package torrents

import (
	"context"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/model"
)

const maxResultsLimit uint = 40

func (s *service) Search(ctx context.Context, query string, hint SearchTypeHint, limit uint) ([]model.Torrent, error) {
	if limit == 0 || limit > maxResultsLimit {
		limit = maxResultsLimit
	}

	found, err := s.provider.SearchTorrents(ctx, query, limit)
	if err != nil {
		return nil, err
	}

	// если кто-то накосячил из провайдеров - исправляем
	if uint(len(found)) > limit {
		found = found[:limit]
	}
	// TODO: менеджмент ссылок
	return found, nil
}
