package server

import (
	"errors"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/model"
	"net/http"

	"github.com/RacoonMediaServer/rms-media-discovery/internal/server/models"
	"github.com/RacoonMediaServer/rms-media-discovery/internal/server/restapi/operations/accounts"
	accounts_service "github.com/RacoonMediaServer/rms-media-discovery/internal/service/accounts"
	"github.com/go-openapi/runtime/middleware"
)

func (s *Server) getAccounts(params accounts.GetAccountsParams, key *models.Principal) middleware.Responder {
	l := s.log.WithField("req", "getAccounts").WithField("key", key.Token)
	l.Debug("Request")
	if !key.Admin {
		l.Warn("Forbidden. Required admin privileges")
		return middleware.Error(http.StatusForbidden, "Forbidden")
	}

	registered, err := s.Accounts.GetAccounts()
	if err != nil {
		l.Errorf("Get accounts failed: %s", err)
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

	l.Debugf("Got %d results", len(payload.Results))

	return accounts.NewGetAccountsOK().WithPayload(&payload)
}

func (s *Server) createAccount(params accounts.CreateAccountParams, key *models.Principal) middleware.Responder {
	l := s.log.WithField("req", "createAccount").WithField("key", key.Token).WithField("account", *params.Account)
	l.Debug("Request")
	if !key.Admin {
		l.Warn("Forbidden. Required admin privileges")
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
		l.Errorf("Request failed: %s", err)
		return accounts.NewCreateAccountInternalServerError()
	}

	return accounts.NewCreateAccountOK().WithPayload(&accounts.CreateAccountOKBody{ID: &acc.Id})
}

func (s *Server) deleteAccount(params accounts.DeleteAccountParams, key *models.Principal) middleware.Responder {
	l := s.log.WithField("req", "deleteAccount").WithField("key", key.Token).WithField("id", params.ID)
	l.Debug("Request")
	if !key.Admin {
		l.Warn("Forbidden. Required admin privileges")
		return middleware.Error(http.StatusForbidden, "Forbidden")
	}

	if err := s.Accounts.DeleteAccount(params.ID); err != nil {
		l.Errorf("Request failed: %s", err)
		if errors.Is(err, accounts_service.ErrNotFound) {
			return accounts.NewDeleteAccountNotFound()
		}
		return accounts.NewDeleteAccountInternalServerError()
	}

	return accounts.NewDeleteAccountOK()
}
