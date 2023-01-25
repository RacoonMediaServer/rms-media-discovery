package torrents

import (
	"context"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/model"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/utils"
)

const maxResultsLimit uint = 40

func (s *service) Search(ctx context.Context, query string, hint SearchTypeHint, limit uint) ([]model.Torrent, error) {
	if limit == 0 || limit > maxResultsLimit {
		limit = maxResultsLimit
	}

	// чистим протухшие ссылки
	s.cleanExpiredLinks()

	found, err := s.provider.SearchTorrents(ctx, query, limit)
	if err != nil {
		return nil, err
	}

	// если кто-то накосячил из провайдеров - исправляем
	found = utils.Bound(found, limit)

	// генерируем ссылки на скачивание
	for i := range found {
		s.processTorrentLink(&found[i])
	}

	return found, nil
}
