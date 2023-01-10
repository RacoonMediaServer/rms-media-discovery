package provider

import (
	"context"

	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/model"
)

// MovieInfoProvider интерфейс сущностей, которые позволяют получать информацию о фильмах и сериалах
type MovieInfoProvider interface {
	SearchMovies(ctx context.Context, query string) ([]model.Movie, error)
}

// TorrentsProvider интерфейс сущностей, которые умеют искать по торрентами
type TorrentsProvider interface {
	SearchTorrents(ctx context.Context, query string) ([]model.Torrent, error)
	Download(ctx context.Context, link string) ([]byte, error)
}
