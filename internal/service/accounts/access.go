package accounts

import (
	"git.rms.local/RacoonMediaServer/rms-media-discovery/pkg/model"
)

func (s *service) GetCredentials(serviceId string) (model.Credentials, error) {
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

	login, ok := acc.Credentials["login"]
	if !ok {
		return model.Credentials{}, ErrNotFound
	}

	password, ok := acc.Credentials["password"]
	if !ok {
		return model.Credentials{}, ErrNotFound
	}

	return model.Credentials{
		AccountId: acc.Id,
		Login:     login,
		Password:  password,
	}, nil
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

	token, ok := acc.Credentials["token"]
	if !ok {
		return model.ApiKey{}, ErrNotFound
	}

	return model.ApiKey{
		AccountId: acc.Id,
		Key:       token,
	}, nil
}

func (s *service) MarkUnaccesible(accountId string) {
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
