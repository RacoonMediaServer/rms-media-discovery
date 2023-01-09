package server

import (
	"net/http"

	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/server/models"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/server/restapi/operations/accounts"
	"github.com/go-openapi/runtime/middleware"
)

func (s *Server) getAccounts(params accounts.GetAccountsParams, key *models.Principal) middleware.Responder {
	if !key.Admin {
		return middleware.Error(http.StatusForbidden, "Forbidden")
	}

	return accounts.NewGetAccountsInternalServerError()
}

func (s *Server) createAccount(params accounts.CreateAccountParams, key *models.Principal) middleware.Responder {
	if !key.Admin {
		return middleware.Error(http.StatusForbidden, "Forbidden")
	}

	return accounts.NewCreateAccountInternalServerError()
}

func (s *Server) deleteAccount(params accounts.DeleteAccountParams, key *models.Principal) middleware.Responder {
	if !key.Admin {
		return middleware.Error(http.StatusForbidden, "Forbidden")
	}

	return accounts.NewDeleteAccountInternalServerError()
}
