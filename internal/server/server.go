package server

import (
	"fmt"

	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/server/restapi"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/server/restapi/operations"
	"github.com/apex/log"
	"github.com/go-openapi/loads"
)

type Server struct {
	srv *restapi.Server
}

func (s *Server) ListenAndServer(host string, port int) error {
	if s.srv == nil {
		swaggerSpec, err := loads.Embedded(restapi.SwaggerJSON, restapi.FlatSwaggerJSON)
		if err != nil {
			return err
		}

		// создаем хендлеры API по умолчанию
		api := operations.NewServerAPI(swaggerSpec)
		// if s.Calendar != nil {
		// 	s.configureAPI(api)
		// }

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
