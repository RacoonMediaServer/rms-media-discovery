package provider

import (
	"context"
	model2 "git.rms.local/RacoonMediaServer/rms-media-discovery/pkg/model"
)

const resultsLimit = 10

type Provider interface {
	ID() string
}

// MovieInfoProvider интерфейс сущностей, которые позволяют получать информацию о фильмах и сериалах
type MovieInfoProvider interface {
	Provider
	SearchMovies(ctx context.Context, query string, limit uint) ([]model2.Movie, error)
}

// TorrentsProvider интерфейс сущностей, которые умеют искать по торрентами
type TorrentsProvider interface {
	Provider
	SearchTorrents(ctx context.Context, query model2.SearchQuery) ([]model2.Torrent, error)
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
