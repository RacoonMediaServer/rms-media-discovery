package torrents

import (
	"context"
	"errors"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/heuristic"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/model"
	"time"
)

const maxResultsLimit uint = 40

func (s *Service) Search(ctx context.Context, q model.SearchQuery) ([]model.Torrent, error) {
	if q.Limit == 0 || q.Limit > maxResultsLimit {
		q.Limit = maxResultsLimit
	}

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

func (s *Service) SearchAsync(query model.SearchQuery) (taskID string, err error) {
	taskID, err = s.gen.Generate()
	if err != nil {
		return
	}

	task := &searchTask{
		q:         query,
		f:         s.Search,
		startTime: time.Now(),
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

func (s *Service) findTask(id string) (*searchTask, error) {
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
func (s *Service) Cancel(taskID string) error {
	task, err := s.findTask(taskID)
	if err != nil {
		return err
	}

	task.cancel()
	s.tasks.Delete(taskID)
	return nil
}

func (s *Service) Status(taskID string) (TaskStatus, error) {
	task, err := s.findTask(taskID)
	if err != nil {
		return TaskStatus{}, err
	}
	st := task.status()
	if st.Status != model.Working {
		s.tasks.Delete(taskID)
	}
	return st, nil
}

func (s *Service) cleanExpiredTask() {
	now := time.Now()

	tmp := map[any]struct{}{}
	s.tasks.Range(func(key, value any) bool {
		task, ok := value.(*searchTask)
		if ok && task.isExpired(now) {
			tmp[key] = struct{}{}
		}
		return true
	})
	for k, _ := range tmp {
		s.log.Debugf("Task '%s' is expired", k.(string))
		s.tasks.Delete(k)
	}
}
