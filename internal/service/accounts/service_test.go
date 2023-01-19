package accounts

import (
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/mocks"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/model"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

var imdb1 = model.Account{
	Id: "imdb.1",
	Credentials: map[string]string{
		"token": "t1",
	},
}

var imdb2 = model.Account{
	Id: "imdb.2",
	Credentials: map[string]string{
		"token": "t2",
	},
}

var myapi1 = model.Account{
	Id: "myapi.1",
	Credentials: map[string]string{
		"login":    "user",
		"password": "pwd",
	},
}

func checkInitialAccounts(t *testing.T, s Service) {
	cred, err := s.GetCredentials("myapi")
	assert.NoError(t, err)
	assert.Equal(t, model.Credentials{
		AccountId: "myapi.1",
		Login:     "user",
		Password:  "pwd",
	}, cred)

	key, err := s.GetApiKey("imdb")
	assert.NoError(t, err)
	assert.Equal(t, model.ApiKey{
		AccountId: "imdb.1",
		Key:       "t1",
	}, key)

	key, err = s.GetApiKey("imdb")
	assert.NoError(t, err)
	assert.Equal(t, model.ApiKey{
		AccountId: "imdb.2",
		Key:       "t2",
	}, key)
}

func newTestService(t *testing.T, m *mocks.MockAccountDatabase) Service {
	s := New(m)

	accounts := []model.Account{
		imdb1, imdb2, myapi1,
	}

	m.EXPECT().LoadAccounts().Return(accounts, nil)
	err := s.Initialize()
	assert.NoError(t, err)

	return s
}

func TestService_Initialize(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mocks.NewMockAccountDatabase(ctrl)
	s := New(m)

	m.EXPECT().LoadAccounts().Return(nil, io.EOF)
	err := s.Initialize()
	assert.ErrorIs(t, err, io.EOF)

	accounts := []model.Account{
		imdb1, imdb2, myapi1,
	}

	m.EXPECT().LoadAccounts().Return(accounts, nil)
	err = s.Initialize()
	assert.NoError(t, err)

	checkInitialAccounts(t, s)
}

func TestService_CreateAccount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mocks.NewMockAccountDatabase(ctrl)
	s := New(m)

	m.EXPECT().LoadAccounts().Return([]model.Account{}, nil)
	err := s.Initialize()
	assert.NoError(t, err)

	err = s.CreateAccount(model.Account{Id: "imdb"})
	assert.Error(t, err)

	m.EXPECT().CreateAccount(gomock.Eq(imdb1)).Return(io.EOF)
	err = s.CreateAccount(imdb1)
	assert.ErrorIs(t, err, io.EOF)

	m.EXPECT().CreateAccount(gomock.Eq(imdb1)).Return(nil)
	err = s.CreateAccount(imdb1)
	assert.NoError(t, err)

	m.EXPECT().CreateAccount(gomock.Eq(imdb2)).Return(nil)
	err = s.CreateAccount(imdb2)
	assert.NoError(t, err)

	m.EXPECT().CreateAccount(gomock.Eq(myapi1)).Return(nil)
	err = s.CreateAccount(myapi1)
	assert.NoError(t, err)

	checkInitialAccounts(t, s)
}

func TestService_DeleteAccount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mocks.NewMockAccountDatabase(ctrl)
	s := newTestService(t, m)

	m.EXPECT().DeleteAccount(gomock.Eq("imdb.5")).Return(io.EOF)
	err := s.DeleteAccount("imdb.5")
	assert.ErrorIs(t, err, io.EOF)

	m.EXPECT().DeleteAccount(gomock.Eq("imdb.5")).Return(nil)
	err = s.DeleteAccount("imdb.5")
	assert.ErrorIs(t, err, ErrNotFound)

	m.EXPECT().DeleteAccount(gomock.Eq("imdb.1")).Return(nil)
	err = s.DeleteAccount("imdb.1")
	assert.NoError(t, err)

	key, err := s.GetApiKey("imdb")
	assert.NoError(t, err)
	assert.Equal(t, model.ApiKey{
		AccountId: "imdb.2",
		Key:       "t2",
	}, key)

	key, err = s.GetApiKey("imdb")
	assert.NoError(t, err)
	assert.Equal(t, model.ApiKey{
		AccountId: "imdb.2",
		Key:       "t2",
	}, key)

	m.EXPECT().DeleteAccount(gomock.Eq("imdb.2")).Return(nil)
	err = s.DeleteAccount("imdb.2")
	assert.NoError(t, err)

	_, err = s.GetApiKey("imdb")
	assert.ErrorIs(t, err, ErrNotFound)

	accounts, err := s.GetAccounts()
	assert.NoError(t, err)
	assert.Equal(t, []model.Account{myapi1}, accounts)

	m.EXPECT().DeleteAccount(gomock.Eq("myapi.1")).Return(nil)
	err = s.DeleteAccount("myapi.1")
	assert.NoError(t, err)

	accounts, err = s.GetAccounts()
	assert.NoError(t, err)
	assert.Equal(t, []model.Account{}, accounts)

	_, err = s.GetApiKey("imdb")
	assert.ErrorIs(t, err, ErrNotFound)

	m.EXPECT().CreateAccount(gomock.Eq(imdb1)).Return(nil)
	err = s.CreateAccount(imdb1)
	assert.NoError(t, err)

	accounts, err = s.GetAccounts()
	assert.NoError(t, err)
	assert.Equal(t, []model.Account{imdb1}, accounts)

	key, err = s.GetApiKey("imdb")
	assert.NoError(t, err)
	assert.Equal(t, model.ApiKey{
		AccountId: "imdb.1",
		Key:       "t1",
	}, key)
}

func TestService_GetAccounts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mocks.NewMockAccountDatabase(ctrl)
	s := newTestService(t, m)

	expected := []model.Account{imdb1, imdb2, myapi1}
	accounts, err := s.GetAccounts()
	assert.NoError(t, err)
	assert.Equal(t, len(expected), len(accounts))

	find := func(accounts []model.Account, id string) *model.Account {
		for i := range accounts {
			if accounts[i].Id == id {
				return &accounts[i]
			}
		}
		return nil
	}

	for _, acc := range expected {
		assert.Equal(t, &acc, find(accounts, acc.Id))
	}
}

func TestService_GetApiKey(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mocks.NewMockAccountDatabase(ctrl)
	s := newTestService(t, m)

	key, err := s.GetApiKey("kinopoisk")
	assert.ErrorIs(t, err, ErrNotFound)

	key, err = s.GetApiKey("imdb")
	assert.NoError(t, err)
	assert.Equal(t, model.ApiKey{
		AccountId: "imdb.1",
		Key:       "t1",
	}, key)

	key, err = s.GetApiKey("imdb")
	assert.NoError(t, err)
	assert.Equal(t, model.ApiKey{
		AccountId: "imdb.2",
		Key:       "t2",
	}, key)

	key, err = s.GetApiKey("imdb")
	assert.NoError(t, err)
	assert.Equal(t, model.ApiKey{
		AccountId: "imdb.1",
		Key:       "t1",
	}, key)

	key, err = s.GetApiKey("myapi")
	assert.ErrorIs(t, err, ErrNotFound)
}

func TestService_GetCredentials(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mocks.NewMockAccountDatabase(ctrl)
	s := newTestService(t, m)

	cred, err := s.GetCredentials("kinopoisk")
	assert.ErrorIs(t, err, ErrNotFound)

	cred, err = s.GetCredentials("imdb")
	assert.ErrorIs(t, err, ErrNotFound)

	cred, err = s.GetCredentials("myapi")
	assert.Equal(t, model.Credentials{
		AccountId: "myapi.1",
		Login:     "user",
		Password:  "pwd",
	}, cred)
}

func TestService_MarkUnaccesible(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mocks.NewMockAccountDatabase(ctrl)
	s := newTestService(t, m)

	cred, err := s.GetCredentials("myapi")
	assert.NoError(t, err)
	assert.Equal(t, model.Credentials{
		AccountId: "myapi.1",
		Login:     "user",
		Password:  "pwd",
	}, cred)
	s.MarkUnaccesible(cred.AccountId)

	cred, err = s.GetCredentials("myapi")
	assert.ErrorIs(t, err, ErrNotFound)

	s.MarkUnaccesible("imdb.1")

	key, err := s.GetApiKey("imdb")
	assert.NoError(t, err)
	assert.Equal(t, model.ApiKey{
		AccountId: "imdb.2",
		Key:       "t2",
	}, key)

	key, err = s.GetApiKey("imdb")
	assert.NoError(t, err)
	assert.Equal(t, model.ApiKey{
		AccountId: "imdb.2",
		Key:       "t2",
	}, key)

	s.MarkUnaccesible("imdb.2")
	key, err = s.GetApiKey("imdb")
	assert.ErrorIs(t, err, ErrNotFound)
}
