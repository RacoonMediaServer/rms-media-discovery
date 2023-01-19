package server

import (
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/model"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/server/models"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/server/restapi/operations/movies"
	"github.com/go-openapi/runtime/middleware"
)

func convertSearchMoviesResult(mov *model.Movie) *models.SearchMoviesResult {
	return &models.SearchMoviesResult{
		Description: mov.Description,
		Genres:      mov.Genres,
		ID:          &mov.ID,
		Poster:      mov.Poster,
		Preview:     mov.Preview,
		Rating:      float64(mov.Rating),
		Seasons:     int64(mov.Seasons),
		Title:       &mov.Title,
		Type:        string(mov.Type),
		Year:        int64(mov.Year),
	}
}

func (s *Server) searchMovies(params movies.SearchMoviesParams, key *models.Principal) middleware.Responder {
	var limit uint
	if params.Limit != nil {
		limit = uint(*params.Limit)
	}
	result, err := s.Movies.Search(params.HTTPRequest.Context(), params.Q, limit)
	if err != nil {
		s.log.Errorf("search movies failed: %s", err)
		return movies.NewSearchMoviesInternalServerError()
	}

	payload := movies.SearchMoviesOKBody{
		Results: make([]*models.SearchMoviesResult, len(result)),
	}

	for i := range result {
		payload.Results[i] = convertSearchMoviesResult(&result[i])
	}

	return movies.NewSearchMoviesOK().WithPayload(&payload)
}
