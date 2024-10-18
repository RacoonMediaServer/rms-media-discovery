package navigator

import (
	"context"
	"net/http"
	"path"
	"time"

	"github.com/apex/log"
	"github.com/playwright-community/playwright-go"
)

const defaultTimeout float64 = 120000.0 * 10

var (
	initialized bool
	environment *playwright.Playwright
	browser     playwright.Browser
	settings    Settings
)

type navigator struct {
	ctx      playwright.BrowserContext
	dumpPath string
	id       string
}

type Navigator interface {
	NewPage(log *log.Entry, ctx context.Context) (Page, error)
	GetCookies(urls ...string) (result []*http.Cookie, err error)
	Close()
}

func Initialize() error {
	var err error

	opts := &playwright.RunOptions{Browsers: []string{"chromium"}}
	if err = playwright.Install(opts); err != nil {
		return err
	}

	if environment, err = playwright.Run(); err != nil {
		return err
	}

	if browser, err = environment.Chromium.Launch(); err != nil {
		return err
	}

	initialized = true

	return nil
}

func New(id string) (Navigator, error) {
	dp := path.Join(settings.DefaultDumpLocation, id)

	ctx, err := browser.NewContext(playwright.BrowserNewContextOptions{Locale: playwright.String("ru-RU")})
	if err != nil {
		return nil, err
	}

	ctx.SetDefaultTimeout(defaultTimeout)
	ctx.SetDefaultNavigationTimeout(defaultTimeout)

	return &navigator{ctx: ctx, dumpPath: dp, id: id}, nil
}

func (n *navigator) NewPage(log *log.Entry, ctx context.Context) (Page, error) {
	p, err := n.ctx.NewPage()
	if err != nil {
		return nil, err
	}

	result := &page{
		ch:       make(chan error),
		ctx:      ctx,
		page:     p,
		log:      log,
		dumpPath: n.dumpPath,
		id:       n.id,
	}

	p.SetDefaultTimeout(defaultTimeout)
	p.SetDefaultNavigationTimeout(defaultTimeout)

	return result, nil
}

func (n *navigator) GetCookies(urls ...string) (result []*http.Cookie, err error) {
	cookies, err := n.ctx.Cookies(urls...)
	if err != nil {
		return
	}

	for _, c := range cookies {
		r := &http.Cookie{
			Name:     c.Name,
			Value:    c.Value,
			Path:     c.Path,
			Domain:   c.Domain,
			Expires:  time.Unix(int64(c.Expires), 0),
			MaxAge:   0,
			Secure:   c.Secure,
			HttpOnly: c.HttpOnly,
		}
		result = append(result, r)
	}
	return
}

func (n *navigator) Close() {
	_ = n.ctx.Close()
}
