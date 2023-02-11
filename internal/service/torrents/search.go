package torrents

import (
	"context"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/pkg/heuristic"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/pkg/model"
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

	// проводим эвристический анализ - определяем дополнительную инфу о раздаче на основе заголовка
	for i := range found {
		found[i].Info = heuristic.ParseTitle(found[i].Title)
		//s.log.Infof("'%s': %+v", found[i].Title, found[i].Info)
	}

	// ранжируем и сортируем результаты
	found = rank(found, q)

	// генерируем ссылки на скачивание
	for i := range found {
		s.processTorrentLink(&found[i])
	}

	return found, nil
}
