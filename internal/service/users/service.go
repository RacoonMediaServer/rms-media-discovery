package users

import (
	"errors"
	"fmt"
	"sync"

	"github.com/RacoonMediaServer/rms-media-discovery/internal/db"
	"github.com/RacoonMediaServer/rms-media-discovery/internal/model"
	"github.com/apex/log"
	uuid "github.com/satori/go.uuid"
)

var ErrUserNotFound = errors.New("user not found")

type Service interface {
	Initialize() error
	CheckAccess(token string) (valid bool, admin bool)

	GetUsers() ([]model.User, error)
	CreateUser(info string, isAdmin bool) (string, error)
	DeleteUser(user string) error
}

type service struct {
	db  db.UserDatabase
	log *log.Entry

	mu    sync.RWMutex
	users map[string]*model.User
}

func New(db db.UserDatabase) Service {
	return &service{
		db:    db,
		log:   log.WithField("from", "users"),
		users: make(map[string]*model.User),
	}
}

func (s *service) Initialize() error {
	users, err := s.db.LoadUsers()
	if err != nil {
		return fmt.Errorf("cannot load users: %+w", err)
	}

	if len(users) == 0 {
		admin := model.User{
			Id:      newUserId(),
			Info:    "Default admin user",
			IsAdmin: true,
		}
		if err := s.db.CreateUser(admin); err != nil {
			return fmt.Errorf("cannot create default admin user: %+w", err)
		}
		s.log.Infof("Default admin user created: %s", admin.Id)
		users = append(users, admin)
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range users {
		s.users[users[i].Id] = &users[i]
	}

	return nil
}

func newUserId() string {
	return uuid.NewV4().String()
}
