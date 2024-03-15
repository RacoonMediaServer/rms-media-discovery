package server

import (
	"context"
	"github.com/RacoonMediaServer/rms-media-discovery/internal/service/torrents"
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

type MusicService interface {
	Search(ctx context.Context, query string, limit uint, searchType model.MusicSearchType) ([]model.Music, error)
}

type TorrentService interface {
	Search(ctx context.Context, query model.SearchQuery) ([]model.Torrent, error)
	SearchAsync(query model.SearchQuery) (taskID string, err error)
	Status(taskID string) (torrents.TaskStatus, error)
	Cancel(taskID string) error
	Download(ctx context.Context, link string) ([]byte, error)
}
