package imdb

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
	"strconv"
	"strings"

	"github.com/apex/log"
)

type imdbProvider struct {
	access model.AccessProvider
	p      pipeline.Pipeline
	r      requester.Requester
}

const (
	imdbEndpoint = "https://imdb-api.com/ru/API"
)

var (
	errBadAccount = errors.New("account is unaccessible")
)

type baseResponse struct {
	ErrorMessage string
}

type searchResponse struct {
	baseResponse
	Results []struct {
		Id         string
		ResultType string
		Title      string
	}
}

type getResponse struct {
	baseResponse
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

func convertInfo(id string, info *getResponse) model.Movie {
	m := model.Movie{
		ID:          id,
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

	return m
}

func NewProvider(access model.AccessProvider) provider.MovieInfoProvider {
	p := &imdbProvider{
		access: access,
		p:      pipeline.Open(pipeline.Settings{Id: "imdb"}),
	}
	p.r = requester.New(p)
	return p
}

func (p *imdbProvider) SearchMovies(ctx context.Context, query string, limit uint) ([]model.Movie, error) {
	l := utils.LogFromContext(ctx, "imdb")
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
		movies = append(movies, convertInfo(item.Id, info))
		if len(movies) >= int(limit) {
			break
		}
	}

	return movies, nil
}

func (p *imdbProvider) search(l *log.Entry, ctx context.Context, query string) (*searchResponse, error) {
	for {
		resp, err := p.p.Do(ctx, func() (interface{}, error) {
			token, err := p.access.GetApiKey("imdb")
			if err != nil {
				return nil, err
			}
			u := fmt.Sprintf("%s/%s/%s/%s", imdbEndpoint, "SearchMovie", token.Key, url.PathEscape(query))
			resp := searchResponse{}
			err = p.r.Get(ctx, u, &resp)

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

		result := resp.(*searchResponse)
		return result, nil
	}
}

func (p *imdbProvider) get(l *log.Entry, ctx context.Context, id string) (*getResponse, error) {
	for {
		resp, err := p.p.Do(ctx, func() (interface{}, error) {
			token, err := p.access.GetApiKey("imdb")
			if err != nil {
				return nil, err
			}
			u := fmt.Sprintf("%s/%s/%s/%s", imdbEndpoint, "Title", token.Key, id)
			resp := getResponse{}
			err = p.r.Get(ctx, u, &resp)

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

		result := resp.(*getResponse)
		return result, nil
	}
}

func (p *imdbProvider) GetMovieInfo(ctx context.Context, id string) (*model.Movie, error) {
	l := utils.LogFromContext(ctx, "imdb")
	info, err := p.get(l, ctx, id)
	if err != nil {
		return nil, err
	}
	m := convertInfo(id, info)
	return &m, nil

}

func (p *imdbProvider) ID() string {
	return "imdb"
}
