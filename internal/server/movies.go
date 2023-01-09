package server

import (
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/server/models"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/server/restapi/operations/movies"
	"github.com/go-openapi/runtime/middleware"
)

func (s *Server) searchMovies(params movies.SearchMoviesParams, key *models.Principal) middleware.Responder {
	return movies.NewSearchMoviesInternalServerError()
}
