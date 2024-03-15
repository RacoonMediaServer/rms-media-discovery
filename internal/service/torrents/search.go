package torrents

import (
	"context"
	"errors"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/heuristic"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/model"
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

func (s *service) SearchAsync(query model.SearchQuery) (taskID string, err error) {
	taskID, err = s.gen.Generate()
	if err != nil {
		return
	}

	task := &searchTask{
		q: query,
		f: s.Search,
	}
	task.ctx, task.cancel = context.WithCancel(s.ctx)
	task.state.ContentType = query.Type
	s.tasks.Store(taskID, task)

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		task.run()
	}()

	return
}

func (s *service) findTask(id string) (*searchTask, error) {
	t, ok := s.tasks.Load(id)
	if !ok {
		return &searchTask{}, ErrTaskNotFound
	}
	task, ok := t.(*searchTask)
	if !ok {
		return &searchTask{}, errors.New("cannot extract searchTask")
	}
	return task, nil
}
func (s *service) Cancel(taskID string) error {
	task, err := s.findTask(taskID)
	if err != nil {
		return err
	}

	task.cancel()
	s.tasks.Delete(taskID)
	return nil
}

func (s *service) Status(taskID string) (TaskStatus, error) {
	task, err := s.findTask(taskID)
	if err != nil {
		return TaskStatus{}, err
	}
	return task.status(), nil
}
