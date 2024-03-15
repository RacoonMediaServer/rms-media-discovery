package accounts

import (
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/model"
)

func (s *Service) GetCredentials(serviceId string) (model.Credentials, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	repo, ok := s.repos[serviceId]
	if !ok {
		return model.Credentials{}, ErrNotFound
	}

	acc, ok := repo.Get()
	if !ok {
		return model.Credentials{}, ErrNotFound
	}

	if acc.Login == "" {
		return model.Credentials{}, ErrNotFound
	}

	if acc.Password == "" {
		return model.Credentials{}, ErrNotFound
	}

	return model.Credentials{
		AccountId: acc.Id,
		Login:     acc.Login,
		Password:  acc.Password,
	}, nil
}

func (s *Service) GetApiKey(serviceId string) (model.ApiKey, error) {
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

	if acc.Token == "" {
		return model.ApiKey{}, ErrNotFound
	}

	return model.ApiKey{
		AccountId: acc.Id,
		Key:       acc.Token,
	}, nil
}

func (s *Service) MarkUnaccesible(accountId string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	m := model.Account{Id: accountId}
	if !m.IsValid() {
		return
	}

	repo, ok := s.repos[m.Service()]
	if !ok {
		return
	}
	repo.MarkUnaccessible(accountId)
}
