package accounts

import (
	"errors"
	"fmt"
	"sync"

	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/db"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/model"
	"github.com/apex/log"
)

var ErrNotFound = errors.New("account not found")

type Service interface {
	model.AccessProvider

	Initialize() error

	GetAccounts() ([]model.Account, error)
	CreateAccount(account model.Account) error
	DeleteAccount(id string) error
}

type service struct {
	db  db.AccountDatabase
	log *log.Entry

	mu    sync.Mutex
	repos map[string]*repository
}

func New(db db.AccountDatabase) Service {
	return &service{
		db:  db,
		log: log.WithField("from", "accounts"),
	}
}

func (s *service) Initialize() error {
	registered, err := s.db.LoadAccounts()
	if err != nil {
		return fmt.Errorf("load accounts from database failed: %+w", err)
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

func (s *service) getOrCreateRepo(serviceId string) *repository {
	repo, ok := s.repos[serviceId]
	if !ok {
		repo = newRepository()
		s.repos[serviceId] = repo
	}

	return repo
}
