package scraper

import (
	"context"
	"github.com/gocolly/colly/v2"
	"net/http"
)

type HTMLCallback func(e *colly.HTMLElement, userData interface{})
type ResponseCallback func(e *colly.Response, userData interface{})

type Selector interface {
	Get(url string) error
	GetAsync(url string, userData interface{})
	Post(url string, data map[string]string) error
	Select(selector string, f HTMLCallback) Selector
	SelectResponse(f ResponseCallback) Selector
}

type Scraper interface {
	SetContext(ctx context.Context)
	Clone() Scraper
	Select(selector string, f HTMLCallback) Selector
	SelectResponse(f ResponseCallback) Selector
	SetCookies(url string, cookies []*http.Cookie) error
	Wait()
}

type scraper struct {
	service string
	c       *colly.Collector
	ctx     context.Context
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
	}
	newScrapper.setCallbacks()
	return newScrapper
}

func New(service string) Scraper {
	s := &scraper{
		service: service,
		c: colly.NewCollector(
			colly.UserAgent("Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36"),
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
	return nil
}

func (s *scraper) Post(url string, data map[string]string) error {
	err := s.c.Post(url, data)
	if err != nil {
		return err
	}

	s.c.Wait()
	return nil
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
	return s.c.SetCookies(url, cookies)
}

func (s *scraper) clear() {
	s.c = s.c.Clone()
	s.setCallbacks()
}
