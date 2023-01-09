package admin

import "git.rms.local/RacoonMediaServer/rms-media-discovery/internal/model"

func (s *service) CheckAccess(token string) (valid bool, admin bool) {
	return
}

func (s *service) GetUsers() ([]model.User, error) {
	return nil, nil
}
func (s *service) CreateUser(info string) (string, error) {
	return "", nil
}
func (s *service) DeleteUser(user string) error {
	return nil
}
