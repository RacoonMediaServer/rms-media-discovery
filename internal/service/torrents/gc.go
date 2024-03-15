package torrents

import "time"

const gcInterval = 10 * time.Second

func (s *Service) gcResourcesProcess() {
	for {
		select {
		case <-time.After(gcInterval):
			s.cleanExpiredLinks()
			s.cleanExpiredTask()
		case <-s.ctx.Done():
			return
		}
	}
}
