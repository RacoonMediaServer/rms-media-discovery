package server

import (
	"fmt"

	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/server/restapi"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/server/restapi/operations"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/service/accounts"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/service/movies"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/service/torrents"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/service/users"
	"github.com/apex/log"
	"github.com/go-openapi/loads"
)

type Server struct {
	srv      *restapi.Server
	Movies   movies.Service
	Torrents torrents.Service
	Users    users.Service
	Accounts accounts.Service
}

func (s *Server) ListenAndServer(host string, port int) error {
	if s.srv == nil {
		swaggerSpec, err := loads.Embedded(restapi.SwaggerJSON, restapi.FlatSwaggerJSON)
		if err != nil {
			return err
		}

		// создаем хендлеры API по умолчанию
		api := operations.NewServerAPI(swaggerSpec)
		s.configureAPI(api)

		// устанавливаем свой логгер
		logCtx := log.WithField("from", "rest")
		api.Logger = func(s string, i ...interface{}) {
			logCtx.Infof(s, i...)
		}

		// создаем и настраиваем сервер
		s.srv = restapi.NewServer(api)
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
