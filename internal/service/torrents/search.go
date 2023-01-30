package torrents

import (
	"context"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/model"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/utils"
)

const maxResultsLimit uint = 40

func (s *service) Search(ctx context.Context, q model.SearchQuery) ([]model.Torrent, error) {
	if q.Limit == 0 || q.Limit > maxResultsLimit {
		q.Limit = maxResultsLimit
	}

	// чистим протухшие ссылки
	s.cleanExpiredLinks()

	found, err := s.provider.SearchTorrents(ctx, q)
	if err != nil {
		return nil, err
	}

	// если кто-то накосячил из провайдеров - исправляем
	found = utils.Bound(found, q.Limit)

	// генерируем ссылки на скачивание
	for i := range found {
		s.processTorrentLink(&found[i])
	}

	return found, nil
}
