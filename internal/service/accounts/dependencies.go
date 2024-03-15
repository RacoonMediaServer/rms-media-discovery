package accounts

import "github.com/RacoonMediaServer/rms-media-discovery/pkg/model"

type AccountDatabase interface {
	LoadAccounts() ([]model.Account, error)
	CreateAccount(account model.Account) error
	DeleteAccount(id string) error
}
