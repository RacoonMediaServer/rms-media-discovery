package provider

import (
	"context"
	"errors"
	"fmt"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/mocks"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/model"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/pipeline"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"os"
	"testing"
)

type overriddenTransportFunc func(req *http.Request) (*http.Response, error)

func (f overriddenTransportFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

type testableProvider interface {
	OverrideTransport(transport http.RoundTripper)
}

type exchange struct {
	verify     func(req *http.Request)
	resp       *http.Response
	mustFailed bool
}

type movieInfoProviderTester struct {
	t  *testing.T
	m  *mocks.MockAccessProvider
	p  MovieInfoProvider
	tp testableProvider
}

func newMovieInfoProviderTester(t *testing.T, m *mocks.MockAccessProvider, p MovieInfoProvider) movieInfoProviderTester {
	underlying, ok := p.(testableProvider)
	assert.True(t, ok)

	return movieInfoProviderTester{
		t:  t,
		m:  m,
		p:  p,
		tp: underlying,
	}
}

func (t *movieInfoProviderTester) testSearchMovieCanceled() {
	t.tp.OverrideTransport(overriddenTransportFunc(func(req *http.Request) (*http.Response, error) {
		return nil, context.Canceled
	}))
	t.m.EXPECT().GetApiKey(gomock.Eq(t.p.ID())).Return(model.ApiKey{
		AccountId: "account",
		Key:       "key",
	}, nil)

	movies, err := t.p.SearchMovies(context.Background(), "Something", 0)
	assert.Nil(t.t, movies)
	assert.Error(t.t, err)
}

func (t *movieInfoProviderTester) testSearchMovieCannotGetApiKey() {
	t.m.EXPECT().GetApiKey(gomock.Eq(t.p.ID())).Return(model.ApiKey{
		AccountId: "account",
		Key:       "key",
	}, nil)
	t.m.EXPECT().MarkUnaccesible(gomock.Eq("account"))
	t.m.EXPECT().GetApiKey(gomock.Eq(t.p.ID())).Return(model.ApiKey{}, errors.New("cannot get api key"))
	t.tp.OverrideTransport(overriddenTransportFunc(func(req *http.Request) (*http.Response, error) {
		return nil, io.EOF
	}))

	movies, err := t.p.SearchMovies(context.Background(), "Something", 0)
	assert.Nil(t.t, movies)
	assert.Error(t.t, err)
}

func (t *movieInfoProviderTester) testSearchMovieMaxAttemptsReached() {
	t.tp.OverrideTransport(overriddenTransportFunc(func(req *http.Request) (*http.Response, error) {
		return nil, io.EOF
	}))

	for i := 0; i < pipeline.MaxAttempts; i++ {
		accountId := fmt.Sprintf("account.%d", i)
		t.m.EXPECT().GetApiKey(gomock.Eq(t.p.ID())).Return(model.ApiKey{
			AccountId: accountId,
			Key:       "key",
		}, nil)
		t.m.EXPECT().MarkUnaccesible(gomock.Eq(accountId))
	}

	movies, err := t.p.SearchMovies(context.Background(), "Something", 0)
	assert.Nil(t.t, movies)
	assert.Error(t.t, err)
}

func (t *movieInfoProviderTester) testSearchMovieUnexpectedStatusCode() {
	t.tp.OverrideTransport(overriddenTransportFunc(func(req *http.Request) (*http.Response, error) {
		resp := http.Response{
			Status:     "Bad request",
			StatusCode: 400,
			Proto:      "HTTP",
			ProtoMajor: 1,
			ProtoMinor: 0,
			Request:    req,
		}
		return &resp, nil
	}))

	for i := 0; i < pipeline.MaxAttempts; i++ {
		accountId := fmt.Sprintf("account.%d", i)
		t.m.EXPECT().GetApiKey(gomock.Eq(t.p.ID())).Return(model.ApiKey{
			AccountId: accountId,
			Key:       "key",
		}, nil)
		t.m.EXPECT().MarkUnaccesible(gomock.Eq(accountId))
	}

	movies, err := t.p.SearchMovies(context.Background(), "Something", 0)
	assert.Nil(t.t, movies)
	assert.Error(t.t, err)
}

func (t *movieInfoProviderTester) testSearchMovie(query string, result []model.Movie, resultErr error, exchanges []exchange) {
	index := 0
	t.tp.OverrideTransport(overriddenTransportFunc(func(req *http.Request) (*http.Response, error) {
		assert.True(t.t, index < len(exchanges), "all exchanges have done")
		exchanges[index].verify(req)
		index++
		return exchanges[index-1].resp, nil
	}))

	for i := range exchanges {
		t.m.EXPECT().GetApiKey(gomock.Eq(t.p.ID())).Return(model.ApiKey{
			AccountId: "account",
			Key:       "key",
		}, nil)
		if exchanges[i].mustFailed {
			t.m.EXPECT().MarkUnaccesible(gomock.Eq("account"))
		}
	}

	movies, err := t.p.SearchMovies(context.Background(), query, 0)
	assert.Equal(t.t, resultErr, err)
	assert.Equal(t.t, result, movies)
}

func (e *exchange) setResponseFile(status int, file string) {
	e.resp = &http.Response{
		Status:     "OK",
		StatusCode: status,
		Proto:      "HTTP",
	}
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	fi, err := f.Stat()
	if err != nil {
		panic(err)
	}
	e.resp.Body = f
	e.resp.ContentLength = fi.Size()
	e.resp.Header = http.Header{}
	e.resp.Header.Add("Content-Type", "application/json")
}

func makeOkExchange(t *testing.T, url, method, contentFile string) exchange {
	exch := exchange{}
	exch.verify = func(req *http.Request) {
		assert.Equal(t, url, req.URL.String())
		assert.Equal(t, method, req.Method)
	}
	exch.setResponseFile(http.StatusOK, contentFile)
	return exch
}
