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
