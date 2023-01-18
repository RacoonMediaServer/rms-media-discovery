package provider

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

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

type imdbListResponse struct {
	Results []struct {
		Id         string
		ResultType string
		Title      string
	}
	ErrorMessage string
}

type imdbResponse struct {
	Title        string
	Image        string
	Type         string
	Year         string
	ErrorMessage string
	Plot         string
	PlotLocal    string
	GenreList    []struct {
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

func (p *imdbProvider) SearchMovies(ctx context.Context, query string) ([]model.Movie, error) {

	list, err := p.search(ctx, query)
	if err != nil {
		return nil, err
	}

	movies := make([]model.Movie, 0)
	for _, item := range list.Results {
		info, err := p.get(ctx, item.Id)
		if err != nil {
			return nil, fmt.Errorf("retrieve info about '%s' failed: %+w", item.Title, err)
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
	}

	return movies, nil
}

func (p *imdbProvider) search(ctx context.Context, query string) (*imdbListResponse, error) {
	resp, err := p.p.Do(ctx, func() pipeline.Result {
		token, err := p.access.GetApiKey("imdb")
		if err != nil {
			return pipeline.Result{Done: true, Err: err}
		}
		u := fmt.Sprintf("%s/%s/%s/%s", imdbEndpoint, "SearchTitle", token.Key, url.PathEscape(query))
		resp := imdbListResponse{}
		err = doRequest(p.cli, ctx, u, &resp)

		if err == nil && resp.ErrorMessage != "" {
			err = fmt.Errorf("imdb response error: %s", resp.ErrorMessage)
		}

		if err != nil {
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				return pipeline.Result{Done: true, Err: err}
			}
			p.access.MarkUnaccesible(token.AccountId)
			return pipeline.Result{Err: err}
		}

		return pipeline.Result{Done: true, Result: &resp}
	})

	if err != nil {
		return nil, err
	}

	result := resp.(*imdbListResponse)
	return result, nil
}

func (p *imdbProvider) get(ctx context.Context, id string) (*imdbResponse, error) {
	resp, err := p.p.Do(ctx, func() pipeline.Result {
		token, err := p.access.GetApiKey("imdb")
		if err != nil {
			return pipeline.Result{Done: true, Err: err}
		}
		u := fmt.Sprintf("%s/%s/%s/%s", imdbEndpoint, "Title", token.Key, id)
		resp := imdbResponse{}
		err = doRequest(p.cli, ctx, u, &resp)

		if err == nil && resp.ErrorMessage != "" {
			err = fmt.Errorf("imdb response error: %s", resp.ErrorMessage)
		}

		if err != nil {
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				return pipeline.Result{Done: true, Err: err}
			}
			p.access.MarkUnaccesible(token.AccountId)
			return pipeline.Result{Err: err}
		}

		return pipeline.Result{Done: true, Result: &resp}
	})

	if err != nil {
		return nil, err
	}

	result := resp.(*imdbResponse)
	return result, nil
}

func (p *imdbProvider) OverrideTransport(transport http.RoundTripper) {
	p.cli.Transport = transport
}

func (p *imdbProvider) ID() string {
	return "imdb"
}
