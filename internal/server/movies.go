package server

import (
	"context"

	"github.com/RacoonMediaServer/rms-media-discovery/internal/server/models"
	"github.com/RacoonMediaServer/rms-media-discovery/internal/server/restapi/operations/movies"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/model"
	"github.com/go-openapi/runtime/middleware"
)

func convertSearchMoviesResult(mov *model.Movie) *models.SearchMoviesResult {
	return &models.SearchMoviesResult{
		Description:   mov.Description,
		Genres:        mov.Genres,
		ID:            &mov.ID,
		Poster:        mov.Poster,
		Preview:       mov.Preview,
		Rating:        float64(mov.Rating),
		Seasons:       int64(mov.Seasons),
		Title:         &mov.Title,
		OriginalTitle: mov.OriginalTitle,
		Type:          string(mov.Type),
		Year:          int64(mov.Year),
	}
}

func (s *Server) searchMovies(params movies.SearchMoviesParams, key *models.Principal) middleware.Responder {
	l := s.log.WithField("key", key.Token).WithField("req", "searchMovies").WithField("q", params.Q)
	l.Debug("Request")

	var limit uint
	if params.Limit != nil {
		limit = uint(*params.Limit)
	}
	result, err := s.Movies.Search(context.WithValue(params.HTTPRequest.Context(), "log", l), params.Q, limit)
	if err != nil {
		l.Errorf("Request failed: %s", err)
		return movies.NewSearchMoviesInternalServerError()
	}

	payload := movies.SearchMoviesOKBody{
		Results: make([]*models.SearchMoviesResult, len(result)),
	}

	for i := range result {
		payload.Results[i] = convertSearchMoviesResult(&result[i])
	}

	l.Debugf("Got %d results", len(result))

	return movies.NewSearchMoviesOK().WithPayload(&payload)
}

func (s *Server) getMovieInfo(params movies.GetMovieInfoParams, key *models.Principal) middleware.Responder {
	l := s.log.WithField("key", key.Token).WithField("req", "getMovieInfo").WithField("id", params.ID)
	l.Debug("Request")

	result, err := s.Movies.Get(context.WithValue(params.HTTPRequest.Context(), "log", l), params.ID)
	if err != nil {
		l.Errorf("Request failed: %s", err)
		return movies.NewGetMovieInfoNotFound()
	}

	if result == nil {
		return movies.NewGetMovieInfoNotFound()
	}

	return movies.NewGetMovieInfoOK().WithPayload(convertSearchMoviesResult(result))
}
