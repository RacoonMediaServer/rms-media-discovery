package admin

import "git.rms.local/RacoonMediaServer/rms-media-discovery/internal/model"

type Service interface {
	GetUsers() ([]model.User, error)
	CreateUser(info string) (string, error)
	DeleteUser(user string) error

	GetAccounts() ([]model.Account, error)
	CreateAccount(account model.Account) (string, error)
	DeleteAccount(id string) error
}

type service struct {
}

func New() Service {
	return &service{}
}
