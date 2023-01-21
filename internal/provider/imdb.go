package provider

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/model"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/pipeline"
	"github.com/apex/log"
)

type imdbProvider struct {
	log    *log.Entry
	access model.AccessProvider
	p      pipeline.Pipeline
	cli    http.Client
}

const imdbEndpoint = "https://imdb-api.com/ru/API"
const resultsLimit = 10

type imdbBaseResponse struct {
	ErrorMessage string
}

type imdbListResponse struct {
	imdbBaseResponse
	Results []struct {
		Id         string
		ResultType string
		Title      string
	}
}

type imdbResponse struct {
	imdbBaseResponse
	Title     string
	Image     string
	Type      string
	Year      string
	Plot      string
	PlotLocal string
	GenreList []struct {
		Key   string
		Value string
	}
	ImDbRating   string
	TvSeriesInfo struct {
		YearEnd string
		Seasons []string
	}
}

func NewImdbProvider(access model.AccessProvider) MovieInfoProvider {
	return &imdbProvider{
		log:    log.WithField("from", "imdb"),
		access: access,
		p:      pipeline.Open(pipeline.Settings{Id: "imdb"}),
	}
}

func (p *imdbProvider) SearchMovies(ctx context.Context, query string, limit uint) ([]model.Movie, error) {

	l := p.log.WithField("query", query)
	l.Info("Searching...")
	list, err := p.search(l, ctx, query)
	if err != nil {
		return nil, err
	}
	l.Infof("Got %d results", len(list.Results))

	movies := make([]model.Movie, 0)
	for _, item := range list.Results {
		info, err := p.get(l, ctx, item.Id)
		if err != nil {
			l.Errorf("Retrieve info about '%s' failed: %s", item.Title, err)
			continue
		}
		m := model.Movie{
			ID:          item.Id,
			Title:       info.Title,
			Description: info.Plot,
			Poster:      info.Image,
			Seasons:     uint(len(info.TvSeriesInfo.Seasons)),
		}

		if info.PlotLocal != "" {
			m.Description = info.PlotLocal
		}

		rating, _ := strconv.ParseFloat(info.ImDbRating, 32)
		m.Rating = float32(rating)

		year, err := strconv.ParseUint(info.Year, 10, 16)
		if err == nil {
			m.Year = uint(year)
		}

		for _, genre := range info.GenreList {
			m.Genres = append(m.Genres, genre.Value)
		}

		m.Type = model.MovieType_Movie
		if info.Type == "TVSeries" {
			m.Type = model.MovieType_TvSeries
		}

		movies = append(movies, m)
		if len(movies) >= resultsLimit {
			break
		}
		if limit != 0 && len(movies) >= int(limit) {
			break
		}
	}

	return movies, nil
}

func (p *imdbProvider) search(l *log.Entry, ctx context.Context, query string) (*imdbListResponse, error) {
	for {
		resp, err := p.p.Do(ctx, func() (interface{}, error) {
			token, err := p.access.GetApiKey("imdb")
			if err != nil {
				return nil, err
			}
			u := fmt.Sprintf("%s/%s/%s/%s", imdbEndpoint, "SearchMovie", token.Key, url.PathEscape(query))
			resp := imdbListResponse{}
			err = doRequest(l, p.cli, ctx, u, &resp)

			if err == nil && resp.ErrorMessage != "" {
				if strings.HasPrefix(resp.ErrorMessage, "Maximum usage") {
					p.access.MarkUnaccesible(token.AccountId)
					return nil, errBadAccount
				}
				err = fmt.Errorf("imdb response error: %s", resp.ErrorMessage)
			}

			if err != nil {
				l.Errorf("Search failed: %s", err)
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

		result := resp.(*imdbListResponse)
		return result, nil
	}
}

func (p *imdbProvider) get(l *log.Entry, ctx context.Context, id string) (*imdbResponse, error) {
	for {
		resp, err := p.p.Do(ctx, func() (interface{}, error) {
			token, err := p.access.GetApiKey("imdb")
			if err != nil {
				return nil, err
			}
			u := fmt.Sprintf("%s/%s/%s/%s", imdbEndpoint, "Title", token.Key, id)
			resp := imdbResponse{}
			err = doRequest(l, p.cli, ctx, u, &resp)

			if err == nil && resp.ErrorMessage != "" {
				if strings.HasPrefix(resp.ErrorMessage, "Maximum usage") {
					p.access.MarkUnaccesible(token.AccountId)
					return nil, errBadAccount
				}
				err = fmt.Errorf("imdb response error: %s", resp.ErrorMessage)
			}

			if err != nil {
				l.Errorf("Get info failed: %s", err)
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

		result := resp.(*imdbResponse)
		return result, nil
	}
}

func (p *imdbProvider) OverrideTransport(transport http.RoundTripper) {
	p.cli.Transport = transport
}

func (p *imdbProvider) ID() string {
	return "imdb"
}
