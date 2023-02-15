package server

import (
	"errors"
	"net/http"

	"github.com/RacoonMediaServer/rms-media-discovery/internal/server/models"
	"github.com/RacoonMediaServer/rms-media-discovery/internal/server/restapi/operations/users"
	usersrv "github.com/RacoonMediaServer/rms-media-discovery/internal/service/users"
	"github.com/go-openapi/runtime/middleware"
)

func (s *Server) getUsers(params users.GetUsersParams, key *models.Principal) middleware.Responder {
	l := s.log.WithField("req", "getUsers").WithField("key", key.Token)
	l.Debug("Request")
	if !key.Admin {
		l.Warn("Forbidden. Required admin privileges")
		return middleware.Error(http.StatusForbidden, "Forbidden")
	}

	registeredUsers, err := s.Users.GetUsers()
	if err != nil {
		s.log.Errorf("Get users failed: %s", err)
		return users.NewGetUsersInternalServerError()
	}

	payload := users.GetUsersOKBody{}
	for _, user := range registeredUsers {
		payload.Results = append(payload.Results, &users.GetUsersOKBodyResultsItems0{
			ID:      user.Id,
			Info:    user.Info,
			IsAdmin: user.IsAdmin,
		})
	}
	return users.NewGetUsersOK().WithPayload(&payload)
}

func (s *Server) createUser(params users.CreateUserParams, key *models.Principal) middleware.Responder {
	l := s.log.WithField("req", "createUser").WithField("key", key.Token)
	if !key.Admin {
		l.Warn("Forbidden. Required admin privileges")
		return middleware.Error(http.StatusForbidden, "Forbidden")
	}

	isAdmin := false
	if params.User.IsAdmin != nil {
		isAdmin = *params.User.IsAdmin
	}
	id, err := s.Users.CreateUser(*params.User.Info, isAdmin)
	if err != nil {
		s.log.Errorf("Create user failed: %s", err)
		return users.NewCreateUserInternalServerError()
	}

	return users.NewCreateUserOK().WithPayload(&users.CreateUserOKBody{ID: &id})
}

func (s *Server) deleteUser(params users.DeleteUserParams, key *models.Principal) middleware.Responder {
	l := s.log.WithField("req", "deleteUser").WithField("key", key.Token).WithField("id", params.ID)
	if !key.Admin {
		l.Warn("Forbidden. Required admin privileges")
		return middleware.Error(http.StatusForbidden, "Forbidden")
	}

	if err := s.Users.DeleteUser(params.ID); err != nil {
		l.Errorf("Request failed: %s", err)
		if errors.Is(err, usersrv.ErrUserNotFound) {
			return users.NewDeleteUserNotFound()
		}
		return users.NewDeleteUserInternalServerError()
	}

	return users.NewDeleteUserOK()
}
