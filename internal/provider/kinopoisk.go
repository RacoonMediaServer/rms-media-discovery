package provider

import (
	"context"

	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/model"
	"github.com/apex/log"
)

type kinopoiskProvider struct {
	log    *log.Entry
	access model.AccessProvider
}

func NewKinopoiskProvider(access model.AccessProvider) MovieInfoProvider {
	return &kinopoiskProvider{
		log:    log.WithField("from", "kinopoisk"),
		access: access,
	}
}

func (p *kinopoiskProvider) SearchMovies(ctx context.Context, query string) ([]model.Movie, error) {
	return nil, nil
}

func (p *kinopoiskProvider) ID() string {
	return "kinopoisk"
}
