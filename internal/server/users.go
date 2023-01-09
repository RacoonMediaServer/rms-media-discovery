package server

import (
	"net/http"

	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/server/models"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/server/restapi/operations/users"
	"github.com/go-openapi/runtime/middleware"
)

func (s *Server) getUsers(params users.GetUsersParams, key *models.Principal) middleware.Responder {
	if !key.Admin {
		return middleware.Error(http.StatusForbidden, "Forbidden")
	}
	return users.NewGetUsersInternalServerError()
}

func (s *Server) createUser(params users.CreateUserParams, key *models.Principal) middleware.Responder {
	if !key.Admin {
		return middleware.Error(http.StatusForbidden, "Forbidden")
	}

	return users.NewCreateUserInternalServerError()
}

func (s *Server) deleteUser(params users.DeleteUserParams, key *models.Principal) middleware.Responder {
	if !key.Admin {
		return middleware.Error(http.StatusForbidden, "Forbidden")
	}

	return users.NewDeleteUserInternalServerError()
}
