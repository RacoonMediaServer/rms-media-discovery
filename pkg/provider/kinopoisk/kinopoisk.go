package kinopoisk

import (
	"context"
	"errors"
	"fmt"
	"github.com/RacoonMediaServer/rms-media-discovery/internal/utils"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/model"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/pipeline"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/provider"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/requester"
	"net/url"
)

type kinopoiskProvider struct {
	access model.AccessProvider
	p      pipeline.Pipeline
	r      requester.Requester
}

const (
	kinopoiskEndpoint = "https://api.kinopoisk.dev/v1.3/movie"
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

func convertInfo(id string, info *getResponse) model.Movie {
	m := model.Movie{
		ID:          id,
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

	return m
}

func NewProvider(access model.AccessProvider) provider.MovieInfoProvider {
	p := &kinopoiskProvider{
		access: access,
		p:      pipeline.Open(pipeline.Settings{Id: "kinopoisk"}),
	}
	p.r = requester.New(p)
	return p
}

func (p *kinopoiskProvider) SearchMovies(ctx context.Context, query string, limit uint) ([]model.Movie, error) {
	l := utils.LogFromContext(ctx, "kinopoisk")
	l.Info("Searching...")
	list, err := p.search(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	l.Infof("Got %d results", len(list.Docs))

	movies := make([]model.Movie, 0)
	for _, item := range list.Docs {
		if len(item.ExternalID.Imdb) == 0 {
			continue
		}
		info, err := p.get(ctx, item.ExternalID.Imdb)
		if err != nil {
			l.Errorf("Retrieve info about '%s' failed: %s", item.Name, err)
			continue
		}
		movies = append(movies, convertInfo(item.ExternalID.Imdb, info))
		if len(movies) >= int(limit) {
			break
		}
	}

	return movies, nil
}

func (p *kinopoiskProvider) ID() string {
	return "kinopoisk"
}

func (p *kinopoiskProvider) search(ctx context.Context, query string, limit uint) (*searchResponse, error) {
	for {
		resp, err := p.p.Do(ctx, func() (interface{}, error) {
			token, err := p.access.GetApiKey("kinopoisk")
			if err != nil {
				return nil, err
			}
			u := fmt.Sprintf("%s/?token=%s&field=names.name&search=%s&limit=%d&sortField=rating.imdb&sortType=-1&isStrict=false", kinopoiskEndpoint, token.Key, url.QueryEscape(query), limit)
			resp := searchResponse{}
			err = p.r.Get(ctx, u, &resp)

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

func (p *kinopoiskProvider) get(ctx context.Context, id string) (*getResponse, error) {
	for {
		resp, err := p.p.Do(ctx, func() (interface{}, error) {
			token, err := p.access.GetApiKey("kinopoisk")
			if err != nil {
				return nil, err
			}
			u := fmt.Sprintf("%s/?token=%s&field=externalId.imdb&search=%s", kinopoiskEndpoint, token.Key, id)
			resp := getResponse{}
			err = p.r.Get(ctx, u, &resp)

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

func (p *kinopoiskProvider) GetMovieInfo(ctx context.Context, id string) (*model.Movie, error) {
	info, err := p.get(ctx, id)
	if err != nil {
		return nil, err
	}
	m := convertInfo(id, info)
	return &m, nil
}
