// Code generated by go-swagger; DO NOT EDIT.

package admin

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
)

// New creates a new admin API client.
func New(transport runtime.ClientTransport, formats strfmt.Registry) ClientService {
	return &Client{transport: transport, formats: formats}
}

/*
Client for admin API
*/
type Client struct {
	transport runtime.ClientTransport
	formats   strfmt.Registry
}

// ClientOption is the option for Client methods
type ClientOption func(*runtime.ClientOperation)

// ClientService is the interface for Client methods
type ClientService interface {
	CreateAccount(params *CreateAccountParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*CreateAccountOK, error)

	CreateUser(params *CreateUserParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*CreateUserOK, error)

	DeleteAccount(params *DeleteAccountParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*DeleteAccountOK, error)

	DeleteUser(params *DeleteUserParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*DeleteUserOK, error)

	GetAccounts(params *GetAccountsParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*GetAccountsOK, error)

	GetUsers(params *GetUsersParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*GetUsersOK, error)

	SetTransport(transport runtime.ClientTransport)
}

/*
CreateAccount создатьs новый аккаунт
*/
func (a *Client) CreateAccount(params *CreateAccountParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*CreateAccountOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewCreateAccountParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "createAccount",
		Method:             "POST",
		PathPattern:        "/admin/accounts",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &CreateAccountReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*CreateAccountOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for createAccount: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
CreateUser создатьs новый ключ пользователя
*/
func (a *Client) CreateUser(params *CreateUserParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*CreateUserOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewCreateUserParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "createUser",
		Method:             "POST",
		PathPattern:        "/admin/users",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &CreateUserReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*CreateUserOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for createUser: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
DeleteAccount удалитьs аккаунт
*/
func (a *Client) DeleteAccount(params *DeleteAccountParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*DeleteAccountOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewDeleteAccountParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "deleteAccount",
		Method:             "DELETE",
		PathPattern:        "/admin/accounts/{id}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &DeleteAccountReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*DeleteAccountOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for deleteAccount: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
DeleteUser удалитьs ключ пользователя
*/
func (a *Client) DeleteUser(params *DeleteUserParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*DeleteUserOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewDeleteUserParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "deleteUser",
		Method:             "DELETE",
		PathPattern:        "/admin/users/{id}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &DeleteUserReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*DeleteUserOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for deleteUser: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
GetAccounts получитьs список список акканутов и токенов к внешним системам
*/
func (a *Client) GetAccounts(params *GetAccountsParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*GetAccountsOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetAccountsParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "getAccounts",
		Method:             "GET",
		PathPattern:        "/admin/accounts",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &GetAccountsReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*GetAccountsOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for getAccounts: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
GetUsers получитьs список пользователей и информацию по ним
*/
func (a *Client) GetUsers(params *GetUsersParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*GetUsersOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetUsersParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "getUsers",
		Method:             "GET",
		PathPattern:        "/admin/users",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &GetUsersReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*GetUsersOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for getUsers: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

// SetTransport changes the transport on the client
func (a *Client) SetTransport(transport runtime.ClientTransport) {
	a.transport = transport
}
