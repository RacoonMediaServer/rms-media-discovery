package provider

import (
	"context"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/model"
)

type Provider interface {
	ID() string
}

// MovieInfoProvider интерфейс сущностей, которые позволяют получать информацию о фильмах и сериалах
type MovieInfoProvider interface {
	Provider
	SearchMovies(ctx context.Context, query string, limit uint) ([]model.Movie, error)
}

// TorrentsProvider интерфейс сущностей, которые умеют искать по торрентами
type TorrentsProvider interface {
	Provider
	SearchTorrents(ctx context.Context, query model.SearchQuery) ([]model.Torrent, error)
}

// Captcha настройки капчи для распознования
type Captcha struct {
	Url           string
	CaseSensitive bool
	MinLength     int
	MaxLength     int
}

// CaptchaSolver интерфейс рапознавателя капчи
type CaptchaSolver interface {
	Provider
	Solve(ctx context.Context, captcha Captcha) (string, error)
}
