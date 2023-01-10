package accounts

import (
	"errors"
	"fmt"

	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/model"
)

func (s *service) GetAccounts() (result []model.Account, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, repo := range s.repos {
		for _, acc := range repo.accounts {
			result = append(result, acc.Account)
		}
	}

	return
}

func (s *service) CreateAccount(account model.Account) error {
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

func (s *service) DeleteAccount(id string) error {
	if err := s.db.DeleteAccount(id); err != nil {
		return fmt.Errorf("delete account from database failed: %+w", err)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	repo := s.getOrCreateRepo(model.Account{Id: id}.Service())
	return repo.Delete(id)
}
