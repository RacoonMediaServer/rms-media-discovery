package accounts

import (
	"errors"
	"fmt"
	"sync"

	"github.com/apex/log"
)

var ErrNotFound = errors.New("account not found")

type Service struct {
	db  AccountDatabase
	log *log.Entry

	mu    sync.Mutex
	repos map[string]*repository
}

func New(db AccountDatabase) *Service {
	return &Service{
		db:  db,
		log: log.WithField("from", "accounts"),
	}
}

func (s *Service) Initialize() error {
	registered, err := s.db.LoadAccounts()
	if err != nil {
		return fmt.Errorf("load accounts from database failed: %w", err)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.repos = map[string]*repository{}

	for _, acc := range registered {
		repo := s.getOrCreateRepo(acc.Service())
		repo.Add(acc)
	}

	return nil
}

func (s *Service) getOrCreateRepo(serviceId string) *repository {
	repo, ok := s.repos[serviceId]
	if !ok {
		repo = newRepository(s.log.WithField("service", serviceId))
		s.repos[serviceId] = repo
	}

	return repo
}
