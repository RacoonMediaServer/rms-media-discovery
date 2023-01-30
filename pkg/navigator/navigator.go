package navigator

import (
	"context"
	"github.com/apex/log"
	"github.com/playwright-community/playwright-go"
	"net/http"
	"time"
)

const defaultTimeout float64 = 4000.0 * 10

var (
	environment *playwright.Playwright
	browser     playwright.Browser
)

type navigator struct {
	ctx      playwright.BrowserContext
	dumpPath string
}

type Navigator interface {
	NewPage(log *log.Entry, ctx context.Context) (Page, error)
	GetCookies(urls ...string) (result []*http.Cookie, err error)
	Close()
}

func init() {
	var err error
	opts := &playwright.RunOptions{SkipInstallBrowsers: true}

	if err = playwright.Install(opts); err != nil {
		panic(err)
	}

	if environment, err = playwright.Run(); err != nil {
		panic(err)
	}

	if browser, err = environment.Chromium.Launch(); err != nil {
		panic(err)
	}
}

func New(dumpPath ...string) (Navigator, error) {
	dp := ""
	if len(dumpPath) > 0 {
		dp = dumpPath[0]
	}
	ctx, err := browser.NewContext(playwright.BrowserNewContextOptions{Locale: playwright.String("ru-RU")})
	if err != nil {
		return nil, err
	}

	return &navigator{ctx: ctx, dumpPath: dp}, nil
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
	}

	p.SetDefaultTimeout(defaultTimeout)

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
