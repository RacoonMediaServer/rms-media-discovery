package accounts

import (
	"errors"
	"fmt"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/model"
)

func (s *Service) GetAccounts() (result []model.Account, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	result = []model.Account{}

	for _, repo := range s.repos {
		for _, acc := range repo.accounts {
			result = append(result, acc.Account)
		}
	}

	return
}

func (s *Service) CreateAccount(account model.Account) error {
	if !account.IsValid() {
		return errors.New("invalid account ID")
	}

	if err := s.db.CreateAccount(account); err != nil {
		return fmt.Errorf("save a new account to database failed: %+w", err)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	repo := s.getOrCreateRepo(account.Service())
	repo.Add(account)

	return nil
}

func (s *Service) DeleteAccount(id string) error {
	if err := s.db.DeleteAccount(id); err != nil {
		return fmt.Errorf("delete account from database failed: %+w", err)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	m := &model.Account{Id: id}
	repo := s.getOrCreateRepo(m.Service())
	return repo.Delete(id)
}
