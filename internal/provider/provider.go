package provider

import (
	"context"
	"errors"

	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/model"
)

var errBadAccount = errors.New("account is unaccessible")

const resultsLimit = 10

// MovieInfoProvider интерфейс сущностей, которые позволяют получать информацию о фильмах и сериалах
type MovieInfoProvider interface {
	ID() string
	SearchMovies(ctx context.Context, query string, limit uint) ([]model.Movie, error)
}

// TorrentsProvider интерфейс сущностей, которые умеют искать по торрентами
type TorrentsProvider interface {
	SearchTorrents(ctx context.Context, query string) ([]model.Torrent, error)
	Download(ctx context.Context, link string) ([]byte, error)
}

type Captcha struct {
	Url           string
	CaseSensitive bool
	MinLength     int
	MaxLength     int
}

type CaptchaSolver interface {
	Solve(ctx context.Context, captcha Captcha) (string, error)
}
