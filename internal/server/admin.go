package server

import (
	"net/http"

	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/server/models"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/server/restapi/operations/admin"
	"github.com/go-openapi/runtime/middleware"
)

func (s *Server) getUsers(params admin.GetUsersParams, key *models.Principal) middleware.Responder {
	if !key.Admin {
		return middleware.Error(http.StatusForbidden, "Forbidden")
	}
	return admin.NewGetUsersInternalServerError()
}

func (s *Server) createUser(params admin.CreateUserParams, key *models.Principal) middleware.Responder {
	if !key.Admin {
		return middleware.Error(http.StatusForbidden, "Forbidden")
	}

	return admin.NewCreateUserInternalServerError()
}

func (s *Server) deleteUser(params admin.DeleteUserParams, key *models.Principal) middleware.Responder {
	if !key.Admin {
		return middleware.Error(http.StatusForbidden, "Forbidden")
	}

	return admin.NewDeleteUserInternalServerError()
}

func (s *Server) getAccounts(params admin.GetAccountsParams, key *models.Principal) middleware.Responder {
	if !key.Admin {
		return middleware.Error(http.StatusForbidden, "Forbidden")
	}

	return admin.NewGetAccountsInternalServerError()
}

func (s *Server) createAccount(params admin.CreateAccountParams, key *models.Principal) middleware.Responder {
	if !key.Admin {
		return middleware.Error(http.StatusForbidden, "Forbidden")
	}

	return admin.NewCreateAccountInternalServerError()
}

func (s *Server) deleteAccount(params admin.DeleteAccountParams, key *models.Principal) middleware.Responder {
	if !key.Admin {
		return middleware.Error(http.StatusForbidden, "Forbidden")
	}

	return admin.NewDeleteAccountInternalServerError()
}
