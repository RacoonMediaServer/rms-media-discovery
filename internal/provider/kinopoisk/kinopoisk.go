package kinopoisk

import (
	"context"
	"errors"
	"fmt"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/model"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/pipeline"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/provider"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/utils"
	"github.com/apex/log"
	"net/http"
	"net/url"
)

type kinopoiskProvider struct {
	log    *log.Entry
	access model.AccessProvider
	p      pipeline.Pipeline
	cli    http.Client
}

const (
	kinopoiskEndpoint = "https://api.kinopoisk.dev/movie"
	resultsLimit      = 10
)

var (
	errBadAccount = errors.New("account is unaccessible")
)

type searchResponse struct {
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

type getResponse struct {
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

func NewKinopoiskProvider(access model.AccessProvider) provider.MovieInfoProvider {
	return &kinopoiskProvider{
		log:    log.WithField("from", "kinopoisk"),
		access: access,
		p:      pipeline.Open(pipeline.Settings{Id: "kinopoisk"}),
	}
}

func (p *kinopoiskProvider) SearchMovies(ctx context.Context, query string, limit uint) ([]model.Movie, error) {
	l := p.log.WithField("query", query)
	if limit == 0 || limit > resultsLimit {
		limit = uint(resultsLimit)
	}
	l.Info("Searching...")
	list, err := p.search(l, ctx, query, limit)
	if err != nil {
		return nil, err
	}
	l.Infof("Got %d results", len(list.Docs))

	movies := make([]model.Movie, 0)
	for _, item := range list.Docs {
		if len(item.ExternalID.Imdb) == 0 {
			continue
		}
		info, err := p.get(l, ctx, item.ExternalID.Imdb)
		if err != nil {
			l.Errorf("Retrieve info about '%s' failed: %s", item.Name, err)
			continue
		}
		m := model.Movie{
			ID:          item.ExternalID.Imdb,
			Title:       info.Name,
			Description: info.Description,
			Poster:      info.Poster.Url,
			Seasons:     uint(len(info.SeasonsInfo)),
			Rating:      info.Rating.Imdb,
			Year:        info.Year,
		}

		for _, genre := range info.Genres {
			m.Genres = append(m.Genres, genre.Name)
		}

		m.Type = model.MovieType_Movie
		if info.Type == "tv-series" {
			m.Type = model.MovieType_TvSeries
		}

		movies = append(movies, m)
		if len(movies) >= int(limit) {
			break
		}
	}

	return movies, nil
}

func (p *kinopoiskProvider) ID() string {
	return "kinopoisk"
}

func (p *kinopoiskProvider) search(l *log.Entry, ctx context.Context, query string, limit uint) (*searchResponse, error) {
	for {
		resp, err := p.p.Do(ctx, func() (interface{}, error) {
			token, err := p.access.GetApiKey("kinopoisk")
			if err != nil {
				return nil, err
			}
			u := fmt.Sprintf("%s/?token=%s&field=names.name&search=%s&limit=%d&sortField=rating.imdb&sortType=-1&isStrict=false", kinopoiskEndpoint, token.Key, url.QueryEscape(query), limit)
			resp := searchResponse{}
			err = utils.Get(l, p.cli, ctx, u, &resp)

			if err != nil {
				return nil, err
			}

			return &resp, nil
		})

		if err != nil {
			if errors.Is(err, errBadAccount) {
				continue
			}

			return nil, err
		}

		result := resp.(*searchResponse)
		return result, nil
	}
}

func (p *kinopoiskProvider) get(l *log.Entry, ctx context.Context, id string) (*getResponse, error) {
	for {
		resp, err := p.p.Do(ctx, func() (interface{}, error) {
			token, err := p.access.GetApiKey("kinopoisk")
			if err != nil {
				return nil, err
			}
			u := fmt.Sprintf("%s/?token=%s&field=externalId.imdb&search=%s", kinopoiskEndpoint, token.Key, id)
			resp := getResponse{}
			err = utils.Get(l, p.cli, ctx, u, &resp)

			if err != nil {
				return nil, err
			}

			return &resp, err
		})

		if err != nil {
			if errors.Is(err, errBadAccount) {
				continue
			}
			return nil, err
		}

		result := resp.(*getResponse)
		return result, nil
	}
}
