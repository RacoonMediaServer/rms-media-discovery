package requester

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testProvider struct{}

func (testProvider) ID() string { return "test" }

func TestStatusError_IsChallenge(t *testing.T) {
	assert.True(t, (&StatusError{Code: 403, Header: http.Header{"Cf-Mitigated": {"challenge"}}}).IsChallenge())
	assert.True(t, (&StatusError{Code: 503, Header: http.Header{"Cf-Mitigated": {"Challenge"}}}).IsChallenge())
	assert.False(t, (&StatusError{Code: 403, Header: http.Header{}}).IsChallenge())
	assert.False(t, (&StatusError{Code: 404, Header: http.Header{"Cf-Mitigated": {"challenge"}}}).IsChallenge())
	assert.False(t, (&StatusError{Code: 403}).IsChallenge())
}

func TestRequester_Download_SendsUserAgentAndReturnsStatusError(t *testing.T) {
	var gotUA string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotUA = r.UserAgent()
		if r.URL.Path == "/challenge" {
			w.Header().Set("cf-mitigated", "challenge")
			w.WriteHeader(http.StatusForbidden)
			return
		}
		_, _ = w.Write([]byte("ok"))
	}))
	defer srv.Close()

	r := NewWithUserAgent(testProvider{}, "my-agent")

	data, _, err := r.Download(context.Background(), srv.URL+"/ok")
	assert.NoError(t, err)
	assert.Equal(t, "ok", string(data))
	assert.Equal(t, "my-agent", gotUA)

	_, _, err = r.Download(context.Background(), srv.URL+"/challenge")
	var statusErr *StatusError
	assert.True(t, errors.As(err, &statusErr))
	assert.Equal(t, 403, statusErr.Code)
	assert.True(t, statusErr.IsChallenge())
}
