package server

import "github.com/RacoonMediaServer/rms-media-discovery/pkg/model"

type AccountsService interface {
	GetAccounts() ([]model.Account, error)
	CreateAccount(account model.Account) error
	DeleteAccount(id string) error
}
