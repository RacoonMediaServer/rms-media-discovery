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

const kinopoiskEndpoint = "https://api.kinopoisk.dev/movie"

type kpListResponse struct {
	Docs []struct {
		Id         uint64
		Name       string
		Type       string
		Year       uint
		ExternalID struct {
			Imdb string
		}
	}
}

type kpResponse struct {
	Id          uint64
	Name        string
	Type        string
	Year        uint
	Description string

	Poster struct {
		Url        string
		PreviewUrl string
	}
	Rating struct {
		Imdb float32
	}
	Genres []struct {
		Name string
	}
	SeasonsInfo []struct {
		Number        uint
		EpisodesCount uint
	}
}

func NewKinopoiskProvider(access model.AccessProvider) MovieInfoProvider {
	return &kinopoiskProvider{
		log:    log.WithField("from", "kinopoisk"),
		access: access,
	}
}

func (p *kinopoiskProvider) SearchMovies(ctx context.Context, query string, limit uint) ([]model.Movie, error) {
	return nil, nil
}

func (p *kinopoiskProvider) ID() string {
	return "kinopoisk"
}
