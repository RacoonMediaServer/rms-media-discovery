package scraper

import (
	"context"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"sync"

	"github.com/RacoonMediaServer/rms-media-discovery/pkg/requester"
	"github.com/gocolly/colly/v2"
)

type HTMLCallback func(e *colly.HTMLElement, userData interface{})
type ResponseCallback func(e *colly.Response, userData interface{})

type Selector interface {
	Get(url string) error
	GetAsync(url string, userData interface{})
	Post(url string, form url.Values) error
	Select(selector string, f HTMLCallback) Selector
	SelectResponse(f ResponseCallback) Selector
}

type Scraper interface {
	SetContext(ctx context.Context)
	Clone() Scraper
	Select(selector string, f HTMLCallback) Selector
	SelectResponse(f ResponseCallback) Selector
	SetCookies(url string, cookies []*http.Cookie) error
	Cookies(url string) []*http.Cookie
	Wait()
}

type scraper struct {
	service string
	c       *colly.Collector
	ctx     context.Context

	// strict - Get/Post возвращают ошибку, если запрос завершился неудачей (по умолчанию такие ошибки только логируются)
	strict bool
	mu     sync.Mutex
	reqErr error
}

func (s *scraper) SetContext(ctx context.Context) {
	s.ctx = ctx
	setCollyContext(s.c, ctx)
}

func (s *scraper) Clone() Scraper {
	newScrapper := &scraper{
		service: s.service,
		c:       s.c.Clone(),
		ctx:     s.ctx,
		strict:  s.strict,
	}
	newScrapper.setCallbacks()
	return newScrapper
}

const defaultUserAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36"

// Options - дополнительные параметры скрапера
type Options struct {
	// UserAgent, по умолчанию - UA обычного десктопного Chrome. Нужно задавать, когда cookies привязаны к определенному UA (cf_clearance)
	UserAgent string

	// FailOnRequestError - Get/Post возвращают ошибку, если сервер ответил не 2xx (*requester.StatusError) или запрос не удался.
	// По умолчанию такие ошибки только логируются, а Get возвращает nil.
	FailOnRequestError bool
}

func New(service string) Scraper {
	return NewWithOptions(service, Options{})
}

// NewWithOptions создает скрапер с заданными параметрами
func NewWithOptions(service string, opts Options) Scraper {
	if opts.UserAgent == "" {
		opts.UserAgent = defaultUserAgent
	}
	s := &scraper{
		service: service,
		strict:  opts.FailOnRequestError,
		c: colly.NewCollector(
			colly.UserAgent(opts.UserAgent),
			colly.AllowURLRevisit(),
			colly.Async(true),
		),
	}
	s.setCallbacks()
	return s
}

func (s *scraper) Get(url string) error {
	err := s.c.Visit(url)
	if err != nil {
		return err
	}

	s.c.Wait()
	return s.takeRequestError()
}

func (s *scraper) Post(url string, form url.Values) error {
	err := s.c.PostRaw(url, []byte(form.Encode()))
	if err != nil {
		return err
	}

	s.c.Wait()
	return s.takeRequestError()
}

// recordRequestError запоминает первую ошибку запроса (вызывается из колбэков colly в разных горутинах)
func (s *scraper) recordRequestError(response *colly.Response, err error) {
	if !s.strict {
		return
	}

	var result error = err
	if response != nil && response.StatusCode != 0 {
		statusErr := &requester.StatusError{Code: response.StatusCode}
		if response.Request != nil && response.Request.URL != nil {
			statusErr.URL = response.Request.URL.String()
		}
		if response.Headers != nil {
			statusErr.Header = *response.Headers
		}
		result = statusErr
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	if s.reqErr == nil {
		s.reqErr = result
	}
}

func (s *scraper) takeRequestError() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	err := s.reqErr
	s.reqErr = nil
	return err
}

func (s *scraper) GetAsync(url string, userData interface{}) {
	ctx := colly.NewContext()
	ctx.Put("userData", userData)
	_ = s.c.Request("GET", url, nil, ctx, nil)
}

func (s *scraper) Select(selector string, f HTMLCallback) Selector {
	s.c.OnHTML(selector, func(element *colly.HTMLElement) {
		userData := element.Request.Ctx.GetAny("userData")
		f(element, userData)
	})

	return s
}

func (s *scraper) SelectResponse(f ResponseCallback) Selector {
	s.c.OnResponse(func(response *colly.Response) {
		userData := response.Request.Ctx.GetAny("userData")
		f(response, userData)
	})

	return s
}

func (s *scraper) Wait() {
	s.c.Wait()
}

func (s *scraper) SetCookies(url string, cookies []*http.Cookie) error {
	j, _ := cookiejar.New(nil)
	s.c.SetCookieJar(j)

	return s.c.SetCookies(url, cookies)
}

func (s *scraper) Cookies(url string) []*http.Cookie {
	return s.c.Cookies(url)
}

func (s *scraper) clear() {
	s.c = s.c.Clone()
	s.setCallbacks()
}
