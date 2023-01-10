package server

import (
	"errors"
	"net/http"

	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/model"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/server/models"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/server/restapi/operations/accounts"
	accounts_service "git.rms.local/RacoonMediaServer/rms-media-discovery/internal/service/accounts"
	"github.com/go-openapi/runtime/middleware"
)

func (s *Server) getAccounts(params accounts.GetAccountsParams, key *models.Principal) middleware.Responder {
	if !key.Admin {
		return middleware.Error(http.StatusForbidden, "Forbidden")
	}

	registered, err := s.Accounts.GetAccounts()
	if err != nil {
		s.log.Errorf("Get accounts failed: %s", err)
		return accounts.NewGetAccountsInternalServerError()
	}

	payload := accounts.GetAccountsOKBody{Results: make([]*models.Account, 0)}

	for _, acc := range registered {
		service := acc.Service()
		payload.Results = append(payload.Results, &models.Account{
			ID:       acc.Id,
			Limit:    int64(acc.Limit),
			Service:  &service,
			Login:    acc.Credentials["login"],
			Password: acc.Credentials["password"],
			Token:    acc.Credentials["token"],
		})
	}

	return accounts.NewGetAccountsOK().WithPayload(&payload)
}

func (s *Server) createAccount(params accounts.CreateAccountParams, key *models.Principal) middleware.Responder {
	if !key.Admin {
		return middleware.Error(http.StatusForbidden, "Forbidden")
	}

	acc := model.Account{
		Limit: uint(params.Account.Limit),
	}

	acc.GenerateId(*params.Account.Service)
	acc.Credentials = map[string]string{}

	if len(params.Account.Token) != 0 {
		acc.Credentials["token"] = params.Account.Token
	}

	if len(params.Account.Login) != 0 {
		acc.Credentials["login"] = params.Account.Login
	}

	if len(params.Account.Password) != 0 {
		acc.Credentials["password"] = params.Account.Password
	}

	if err := s.Accounts.CreateAccount(acc); err != nil {
		s.log.Errorf("Create account failed: %s", err)
		return accounts.NewCreateAccountInternalServerError()
	}

	return accounts.NewCreateAccountOK().WithPayload(&accounts.CreateAccountOKBody{ID: &acc.Id})
}

func (s *Server) deleteAccount(params accounts.DeleteAccountParams, key *models.Principal) middleware.Responder {
	if !key.Admin {
		return middleware.Error(http.StatusForbidden, "Forbidden")
	}

	if err := s.Accounts.DeleteAccount(params.ID); err != nil {
		s.log.Errorf("Delete account '%s' failed: %s", params.ID, err)
		if errors.Is(err, accounts_service.ErrNotFound) {
			return accounts.NewDeleteAccountNotFound()
		}
		return accounts.NewDeleteAccountInternalServerError()
	}

	return accounts.NewDeleteAccountOK()
}
