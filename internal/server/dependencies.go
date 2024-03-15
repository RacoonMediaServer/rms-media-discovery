package server

import (
	"context"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/model"
)

type AccountsService interface {
	GetAccounts() ([]model.Account, error)
	CreateAccount(account model.Account) error
	DeleteAccount(id string) error
}

type MoviesService interface {
	Search(ctx context.Context, query string, limit uint) ([]model.Movie, error)
	Get(ctx context.Context, id string) (*model.Movie, error)
}
