package tmdb

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/RacoonMediaServer/rms-media-discovery/internal/utils"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/model"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/pipeline"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/provider"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/requester"
	"github.com/ryanbradynd05/go-tmdb"
	"golang.org/x/time/rate"
)

const requestsPerSecond = 40

type tmdbProvider struct {
	access model.AccessProvider
	p      pipeline.Pipeline
	r      requester.Requester
	g      genreList
	mirror provider.MirrorService
}

var (
	errBadAccount = errors.New("account is unaccessible")
)

func (p *tmdbProvider) GetMovieInfo(ctx context.Context, id string) (*model.Movie, error) {
	if strings.HasPrefix(id, "tt") {
		return p.getImdbMovieInfo(ctx, id)
	}
	if strings.HasPrefix(id, "tmdb_m_") {
		id = strings.TrimPrefix(id, "tmdb_m_")
		intId, err := strconv.ParseInt(id, 10, 32)
		if err != nil {
			return nil, err
		}
		return p.getTmdbMovieInfo(ctx, int(intId))
	}
	if strings.HasPrefix(id, "tmdb_s_") {
		id = strings.TrimPrefix(id, "tmdb_s_")
		intId, err := strconv.ParseInt(id, 10, 32)
		if err != nil {
			return nil, err
		}
		return p.getTmdbTvSeriesInfo(ctx, int(intId))
	}

	return nil, errors.New("invalid ID")
}

func (p *tmdbProvider) getImdbMovieInfo(ctx context.Context, id string) (*model.Movie, error) {
	resp, err := p.request(ctx, func(api *tmdb.TMDb) (interface{}, error) {
		return api.GetFind(id, "imdb_id", map[string]string{"language": "ru-RU"})
	})
	if err != nil {
		return nil, err
	}
	result := resp.(*tmdb.FindResults)
	if len(result.TvResults) != 0 {
		r := &result.TvResults[0]
		seasons, err := p.getSeasons(ctx, r.ID)
		if err != nil {
			return nil, err
		}
		m := &model.Movie{
			ID:            id,
			Description:   r.Overview,
			Genres:        p.g.get(castGenreList(r.GenreIDs)),
			Poster:        p.composePosterURL(r.PosterPath),
			Preview:       p.composePosterURL(r.PosterPath),
			Rating:        r.VoteAverage,
			Seasons:       seasons,
			Title:         r.Name,
			OriginalTitle: r.OriginalName,
			Type:          model.MovieType_TvSeries,
			Year:          parseYear(r.FirstAirDate),
		}
		return m, nil
	} else if len(result.MovieResults) != 0 {
		r := &result.MovieResults[0]
		m := &model.Movie{
			ID:            id,
			Description:   r.Overview,
			Genres:        p.g.get(castGenreList(r.GenreIDs)),
			Poster:        p.composePosterURL(r.PosterPath),
			Preview:       p.composePosterURL(r.PosterPath),
			Rating:        r.VoteAverage,
			Title:         r.Title,
			OriginalTitle: r.OriginalTitle,
			Type:          model.MovieType_Movie,
			Year:          parseYear(r.ReleaseDate),
		}
		return m, nil
	} else {
		return nil, errors.New("nothing found")
	}
}

func (p *tmdbProvider) getTmdbMovieInfo(ctx context.Context, id int) (*model.Movie, error) {
	resp, err := p.request(ctx, func(api *tmdb.TMDb) (interface{}, error) {
		return api.GetMovieInfo(id, map[string]string{"language": "ru-RU"})
	})
	if err != nil {
		return nil, err
	}
	info := resp.(*tmdb.Movie)
	m := &model.Movie{
		ID:            fmt.Sprintf("tmdb_m_%d", info.ID),
		Description:   info.Overview,
		Poster:        p.composePosterURL(info.PosterPath),
		Preview:       p.composePosterURL(info.PosterPath),
		Rating:        info.VoteAverage,
		Title:         info.Title,
		OriginalTitle: info.OriginalTitle,
		Type:          model.MovieType_Movie,
		Year:          parseYear(info.ReleaseDate),
	}
	for _, g := range info.Genres {
		m.Genres = append(m.Genres, strings.ToLower(g.Name))
	}
	return m, nil
}

func (p *tmdbProvider) getTmdbTvSeriesInfo(ctx context.Context, id int) (*model.Movie, error) {
	resp, err := p.request(ctx, func(api *tmdb.TMDb) (interface{}, error) {
		return api.GetTvInfo(id, map[string]string{"language": "ru-RU"})
	})
	if err != nil {
		return nil, err
	}
	info := resp.(*tmdb.TV)
	m := &model.Movie{
		ID:            fmt.Sprintf("tmdb_s_%d", info.ID),
		Description:   info.Overview,
		Poster:        p.composePosterURL(info.PosterPath),
		Preview:       p.composePosterURL(info.PosterPath),
		Rating:        info.VoteAverage,
		Title:         info.Name,
		OriginalTitle: info.OriginalName,
		Seasons:       uint(info.NumberOfSeasons),
		Type:          model.MovieType_TvSeries,
		Year:          parseYear(info.FirstAirDate),
	}
	for _, g := range info.Genres {
		m.Genres = append(m.Genres, strings.ToLower(g.Name))
	}
	return m, nil
}

func (p *tmdbProvider) composePosterURL(path string) string {
	originURL := fmt.Sprintf("https://image.tmdb.org/t/p/w780%s", path)
	return p.mirror.MakeURL(originURL)

}

func parseYear(text string) uint {
	t, err := time.Parse(time.DateOnly, text)
	if err != nil {
		return 0
	}
	return uint(t.Year())
}

func NewProvider(access model.AccessProvider, mirror provider.MirrorService) provider.MovieInfoProvider {
	settings := pipeline.Settings{
		Id:    "tmdb",
		Limit: rate.NewLimiter(1, requestsPerSecond),
	}
	p := &tmdbProvider{
		access: access,
		p:      pipeline.Open(settings),
		mirror: mirror,
	}
	p.r = requester.New(p)
	return p
}

func (p *tmdbProvider) SearchMovies(ctx context.Context, query string, limit uint) ([]model.Movie, error) {
	l := utils.LogFromContext(ctx, "tmdb")
	l.Info("Searching...")
	p.initGenreList(ctx)

	results, err := p.search(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	l.Infof("Got %d results", len(results))
	return results, nil
}

func (p *tmdbProvider) search(ctx context.Context, query string, limit uint) ([]model.Movie, error) {
	resp, err := p.request(ctx, func(api *tmdb.TMDb) (interface{}, error) {
		return api.SearchMulti(query, map[string]string{"language": "ru-RU"})
	})
	if err != nil {
		return nil, err
	}

	result := resp.(*tmdb.MultiSearchResults)
	var movies []model.Movie
	for _, r := range result.Results {

		switch info := r.(type) {
		case *tmdb.MultiSearchTvInfo:
			m := model.Movie{
				ID:            fmt.Sprintf("tmdb_s_%d", info.ID),
				Description:   info.Overview,
				Genres:        p.g.get(castGenreList(info.GenreIDs)),
				Poster:        p.composePosterURL(info.PosterPath),
				Preview:       p.composePosterURL(info.PosterPath),
				Rating:        info.VoteAverage,
				Title:         info.Name,
				OriginalTitle: info.OriginalName,
				Type:          model.MovieType_TvSeries,
				Year:          parseYear(info.FirstAirDate),
			}
			seasons, err := p.getSeasons(ctx, info.ID)
			if err != nil {
				return nil, fmt.Errorf("get seasons count failed: %w", err)
			}
			m.Seasons = seasons
			movies = append(movies, m)
		case *tmdb.MultiSearchMovieInfo:
			m := model.Movie{
				ID:            fmt.Sprintf("tmdb_m_%d", info.ID),
				Description:   info.Overview,
				Genres:        p.g.get(castGenreList(info.GenreIDs)),
				Poster:        p.composePosterURL(info.PosterPath),
				Preview:       p.composePosterURL(info.PosterPath),
				Rating:        info.VoteAverage,
				Title:         info.Title,
				OriginalTitle: info.OriginalTitle,
				Type:          model.MovieType_Movie,
				Year:          parseYear(info.ReleaseDate),
			}
			movies = append(movies, m)
		}
	}

	return movies, nil

}

func (p *tmdbProvider) getSeasons(ctx context.Context, id int) (uint, error) {
	resp, err := p.request(ctx, func(api *tmdb.TMDb) (interface{}, error) {
		return api.GetTvInfo(id, map[string]string{})
	})
	if err != nil {
		return 0, err
	}
	result := resp.(*tmdb.TV)
	return uint(result.NumberOfSeasons), nil
}

func (p *tmdbProvider) ID() string {
	return "tmdb"
}
