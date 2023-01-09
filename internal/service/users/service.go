package users

import "git.rms.local/RacoonMediaServer/rms-media-discovery/internal/model"

type Service interface {
	CheckAccess(token string) (valid bool, admin bool)

	GetUsers() ([]model.User, error)
	CreateUser(info string) (string, error)
	DeleteUser(user string) error
}

type service struct {
}

func New() Service {
	return &service{}
}
