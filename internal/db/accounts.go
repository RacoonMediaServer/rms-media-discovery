package db

import "git.rms.local/RacoonMediaServer/rms-media-discovery/internal/model"

type AccountDatabase interface {
	LoadAccounts() ([]model.Account, error)
	CreateAccount(account model.Account) error
	DeleteAccount(id string) error
}
