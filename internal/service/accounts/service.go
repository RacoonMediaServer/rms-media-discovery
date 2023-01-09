package accounts

import "git.rms.local/RacoonMediaServer/rms-media-discovery/internal/model"

type Service interface {
	GetAccounts() ([]model.Account, error)
	CreateAccount(account model.Account) (string, error)
	DeleteAccount(id string) error
}

type service struct {
}

func New() Service {
	return &service{}
}
