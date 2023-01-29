package navigator

import (
	"bytes"
	"context"
	"github.com/PuerkitoBio/goquery"
	"github.com/apex/log"
	"github.com/playwright-community/playwright-go"
	"io/ioutil"
	"net/http"
	"time"
)

const defaultTimeout float64 = 4000.0 * 10

var (
	environment *playwright.Playwright
	browser     playwright.Browser
)

type Navigator struct {
	ctx      playwright.BrowserContext
	dumpPath string
}

type Page struct {
	ctx      context.Context
	ch       chan error
	page     playwright.Page
	log      *log.Entry
	err      error
	batch    string
	dumpPath string
	doc      *goquery.Document
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

func New(dumpPath ...string) (*Navigator, error) {
	dp := ""
	if len(dumpPath) > 0 {
		dp = dumpPath[0]
	}
	ctx, err := browser.NewContext(playwright.BrowserNewContextOptions{Locale: playwright.String("ru-RU")})
	if err != nil {
		return nil, err
	}

	return &Navigator{ctx: ctx, dumpPath: dp}, nil
}

func (n *Navigator) NewPage(log *log.Entry, ctx context.Context) (*Page, error) {
	page, err := n.ctx.NewPage()
	if err != nil {
		return nil, err
	}

	result := &Page{
		ch:       make(chan error),
		ctx:      ctx,
		page:     page,
		log:      log,
		dumpPath: n.dumpPath,
	}

	page.SetDefaultTimeout(defaultTimeout)

	return result, nil
}

func (p *Page) Batch(title string) *Page {
	p.batch = title
	return p
}

func (p *Page) Goto(url string) *Page {
	if p.err != nil {
		return nil
	}

	go func() {
		p.goTo(url)
	}()

	select {
	case p.err = <-p.ch:
	case <-p.ctx.Done():
		p.err = p.ctx.Err()
	}
	p.checkError("Goto")
	return p
}

func (p *Page) goTo(url string) {

	p.log.Debugf("%s: navigating to '%s'...", p.batch, url)
	_, err := p.page.Goto(url)
	select {
	case p.ch <- err:
	case <-p.ctx.Done():
	}
}

func (p *Page) Error() error {
	return p.err
}

func (p *Page) RaiseError(err error) {
	p.err = err
	p.checkError("Raise")
}

func (p *Page) ClearError() *Page {
	p.err = nil
	return p
}

func (p *Page) Type(selector, text string) *Page {
	if p.err != nil {
		return p
	}
	p.log.Debugf("%s: typing '%s' to '%s'...", p.batch, text, selector)
	p.err = p.page.Type(selector, text)

	p.checkError("Type")
	return p
}

func (p *Page) Click(selector string) *Page {
	if p.err != nil {
		return p
	}

	go p.click(selector)

	select {
	case p.err = <-p.ch:
	case <-p.ctx.Done():
		p.err = p.ctx.Err()
	}

	p.checkError("Click")
	return p
}

func (p *Page) click(selector string) {
	p.log.Debugf("%s: clicking on '%s'...", p.batch, selector)
	err := p.page.Click(selector)
	select {
	case p.ch <- err:
	case <-p.ctx.Done():
	}
}

func (p *Page) FetchContent() *Page {
	if p.err != nil {
		return p
	}

	p.log.Debugf("%s: fetching content", p.batch)
	output := ""
	output, p.err = p.page.Content()
	if p.err == nil {
		p.doc, p.err = goquery.NewDocumentFromReader(bytes.NewReader([]byte(output)))
	}
	p.checkError("FetchContent")
	return p
}

func (p *Page) Document() *goquery.Document {
	return p.doc
}

func (p *Page) Screenshot(fileName string) *Page {
	data, err := p.page.Screenshot()
	if err != nil {
		p.log.Warnf("Cannot get screenshot: %+v", err)
		return p
	}
	if err = ioutil.WriteFile(fileName, data, 0644); err != nil {
		p.log.Warnf("Cannot save screenshot: %+v", err)
	}
	return p
}

func (p *Page) TracePage(fileName string) *Page {
	content, err := p.page.Content()
	if err != nil {
		p.log.Warnf("Cannot fetch page content: %+v", err)
		return p
	}
	if err = ioutil.WriteFile(fileName, []byte(content), 0644); err != nil {
		p.log.Warnf("Cannot save page content: %+v", err)
	}
	return p
}

func (p *Page) Sleep(d time.Duration) *Page {
	if p.err != nil {
		return p
	}
	p.log.Debugf("%s: waiting for %+v", p.batch, d)
	<-time.After(d)
	return p
}

func (p *Page) Address() string {
	return p.page.URL()
}

func (p *Page) Close() {
	_ = p.page.Close()
}

func (n *Navigator) GetCookies(urls ...string) (result []*http.Cookie, err error) {
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

func (n *Navigator) Close() {
	_ = n.ctx.Close()
}
