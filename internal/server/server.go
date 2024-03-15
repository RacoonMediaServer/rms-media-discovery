package server

import (
	"fmt"
	rms_users "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-users"
	"net/http"

	"github.com/RacoonMediaServer/rms-media-discovery/internal/server/restapi"
	"github.com/RacoonMediaServer/rms-media-discovery/internal/server/restapi/operations"
	"github.com/RacoonMediaServer/rms-media-discovery/internal/service/torrents"
	"github.com/apex/log"
	"github.com/go-openapi/loads"
)

type Server struct {
	srv *restapi.Server
	log *log.Entry

	Movies   MoviesService
	Music    MusicService
	Torrents torrents.Service
	Users    rms_users.RmsUsersService
	Accounts AccountsService
}

type monitor struct {
	handler http.Handler
}

func (s *Server) ListenAndServer(host string, port int) error {
	s.log = log.WithField("from", "rest")

	if s.srv == nil {
		swaggerSpec, err := loads.Embedded(restapi.SwaggerJSON, restapi.FlatSwaggerJSON)
		if err != nil {
			return err
		}

		// создаем хендлеры API по умолчанию
		api := operations.NewServerAPI(swaggerSpec)
		s.configureAPI(api)

		// middleware для для тяжелых запросов
		api.AddMiddlewareFor("GET", "/movies/search", getSearchMiddleware("movies"))
		api.AddMiddlewareFor("GET", "/torrents/search", getSearchMiddleware("torrents"))
		api.AddMiddlewareFor("POST", "/torrents/search:run", getSearchMiddleware("torrents"))

		// устанавливаем свой логгер
		api.Logger = func(content string, i ...interface{}) {
			s.log.Infof(content, i...)
		}

		// создаем и настраиваем сервер
		s.srv = restapi.NewServer(api)

		// устанавливаем middleware
		s.srv.SetHandler(setupGlobalMiddleware(api.Serve(nil)))
	}

	s.srv.Host = host
	s.srv.Port = port

	if err := s.srv.Listen(); err != nil {
		return fmt.Errorf("cannot start server: %w", err)
	}

	return s.srv.Serve()
}

func (s *Server) Shutdown() error {
	if s.srv != nil {
		return s.srv.Shutdown()
	}

	return nil
}
