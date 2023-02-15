package users

import (
	"fmt"

	"github.com/RacoonMediaServer/rms-media-discovery/internal/model"
)

func (s *service) CheckAccess(token string) (valid bool, admin bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, ok := s.users[token]
	if !ok {
		return
	}

	s.log.Debugf("Access for user '%s' granted", token)
	userRequestsCounter.WithLabelValues(token).Inc()

	valid = true
	admin = user.IsAdmin
	return
}

func (s *service) GetUsers() (result []model.User, err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	s.log.Debugf("Users: %+v", s.users)

	for _, user := range s.users {
		result = append(result, *user)
	}

	return
}

func (s *service) CreateUser(info string, isAdmin bool) (string, error) {
	user := model.User{
		Id:      newUserId(),
		Info:    info,
		IsAdmin: isAdmin,
	}
	if err := s.db.CreateUser(user); err != nil {
		return "", fmt.Errorf("cannot store new user: %+w", err)
	}

	s.log.Infof("User '%s' created", user.Id)

	s.mu.Lock()
	defer s.mu.Unlock()

	s.users[user.Id] = &user
	return user.Id, nil
}

func (s *service) DeleteUser(user string) error {
	if err := s.db.DeleteUser(user); err != nil {
		return fmt.Errorf("delete user from database failed: %w", err)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.users[user]; !ok {
		return ErrUserNotFound
	}

	delete(s.users, user)
	s.log.Infof("User '%s' deleted", user)
	return nil
}
