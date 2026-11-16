package accounts

import (
	"testing"

	"github.com/RacoonMediaServer/rms-media-discovery/internal/mocks"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/model"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestService_GetCredentials_Session(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mocks.NewMockAccountDatabase(ctrl)
	s := New(m)

	m.EXPECT().LoadAccounts().Return([]model.Account{
		// только cookies: логин и пароль не нужны
		{Id: "rutracker.1", Cookies: "bb_session=1; cf_clearance=2", UserAgent: "UA"},
		// cookies без User-Agent бесполезны (cf_clearance привязан к UA)
		{Id: "bad.1", Cookies: "bb_session=1"},
		// если заданы и cookies, и логин с паролем, приоритет у cookies
		{Id: "both.1", Login: "user", Password: "pwd", Cookies: "c=1", UserAgent: "UA2"},
	}, nil)
	assert.NoError(t, s.Initialize())

	cred, err := s.GetCredentials("rutracker")
	assert.NoError(t, err)
	assert.Equal(t, model.Credentials{
		AccountId: "rutracker.1",
		Cookies:   "bb_session=1; cf_clearance=2",
		UserAgent: "UA",
	}, cred)

	_, err = s.GetCredentials("bad")
	assert.ErrorIs(t, err, ErrNotFound)

	cred, err = s.GetCredentials("both")
	assert.NoError(t, err)
	assert.Equal(t, "c=1", cred.Cookies)
	assert.Equal(t, "UA2", cred.UserAgent)
	assert.Equal(t, "user", cred.Login)
}
