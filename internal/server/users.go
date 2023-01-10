package server

import (
	"errors"
	"net/http"

	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/server/models"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/server/restapi/operations/users"
	usersrv "git.rms.local/RacoonMediaServer/rms-media-discovery/internal/service/users"
	"github.com/go-openapi/runtime/middleware"
)

func (s *Server) getUsers(params users.GetUsersParams, key *models.Principal) middleware.Responder {
	if !key.Admin {
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
	if !key.Admin {
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
	if !key.Admin {
		return middleware.Error(http.StatusForbidden, "Forbidden")
	}

	if err := s.Users.DeleteUser(params.ID); err != nil {
		s.log.Errorf("Delete user failed: %s", err)
		if errors.Is(err, usersrv.ErrUserNotFound) {
			return users.NewDeleteUserNotFound()
		}
		return users.NewDeleteUserInternalServerError()
	}

	return users.NewDeleteUserOK()
}
