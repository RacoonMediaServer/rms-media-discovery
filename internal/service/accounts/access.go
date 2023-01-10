package accounts

import "git.rms.local/RacoonMediaServer/rms-media-discovery/internal/model"

func (s *service) GetCredentials(serviceId string) (model.Credentials, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	repo, ok := s.repos[serviceId]
	if !ok {
		return model.Credentials{}, ErrNotFound
	}

	return repo.GetCredentials()
}

func (s *service) GetApiKey(serviceId string) (model.ApiKey, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	repo, ok := s.repos[serviceId]
	if !ok {
		return model.ApiKey{}, ErrNotFound
	}

	acc, ok := repo.Get()
	if !ok {
		return model.ApiKey{}, ErrNotFound
	}
}

func (s *service) MarkUnaccesible(accountId string) {
	s.mu.Lock()
	defer s.mu.Unlock()

}
