package navigator

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/pkg/requester"
	"github.com/PuerkitoBio/goquery"
	"github.com/apex/log"
	"github.com/playwright-community/playwright-go"
	"github.com/prometheus/client_golang/prometheus"
	uuid "github.com/satori/go.uuid"
	"os"
	"path"
	"time"
)

type page struct {
	ctx      context.Context
	ch       chan error
	page     playwright.Page
	log      *log.Entry
	err      error
	batch    string
	dumpPath string
	doc      *goquery.Document
	id       string
}

type Page interface {
	Batch(title string) Page
	Goto(url string) Page
	RaiseError(err error) Page
	ClearError() Page
	Type(selector, text string) Page
	Click(selector string) Page
	FetchContent() Page
	Screenshot(fileName string) Page
	TracePage(fileName string) Page
	Sleep(d time.Duration) Page

	Document() *goquery.Document
	Address() string
	Error() error

	Close()
}

func (p *page) Batch(title string) Page {
	p.batch = title
	return p
}

func (p *page) Goto(url string) Page {
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

func (p *page) goTo(url string) {

	p.log.Debugf("%s: navigating to '%s'...", p.batch, url)
	timer := prometheus.NewTimer(requester.OutgoingRequestsMetric.WithLabelValues(p.id))
	resp, err := p.page.Goto(url)
	defer func() {
		timer.ObserveDuration()
		if resp != nil {
			requester.OutgoingRequestsCounter.WithLabelValues(fmt.Sprintf("%d", resp.Status()), p.id).Inc()
		}
	}()

	select {
	case p.ch <- err:
	case <-p.ctx.Done():
	}
}

func (p *page) Error() error {
	return p.err
}

func (p *page) RaiseError(err error) Page {
	p.err = err
	p.checkError("Raise")
	return p
}

func (p *page) ClearError() Page {
	p.err = nil
	return p
}

func (p *page) Type(selector, text string) Page {
	if p.err != nil {
		return p
	}
	p.log.Debugf("%s: typing '%s' to '%s'...", p.batch, text, selector)
	p.err = p.page.Type(selector, text)

	p.checkError("Type")
	return p
}

func (p *page) Click(selector string) Page {
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

func (p *page) click(selector string) {
	p.log.Debugf("%s: clicking on '%s'...", p.batch, selector)
	err := p.page.Click(selector)
	select {
	case p.ch <- err:
	case <-p.ctx.Done():
	}
}

func (p *page) FetchContent() Page {
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

func (p *page) Document() *goquery.Document {
	return p.doc
}

func (p *page) Screenshot(fileName string) Page {
	data, err := p.page.Screenshot()
	if err != nil {
		p.log.Warnf("Cannot get screenshot: %+v", err)
		return p
	}
	if err = os.WriteFile(fileName, data, 0644); err != nil {
		p.log.Warnf("Cannot save screenshot: %+v", err)
	}
	return p
}

func (p *page) TracePage(fileName string) Page {
	content, err := p.page.Content()
	if err != nil {
		p.log.Warnf("Cannot fetch page content: %+v", err)
		return p
	}
	if err = os.WriteFile(fileName, []byte(content), 0644); err != nil {
		p.log.Warnf("Cannot save page content: %+v", err)
	}
	return p
}

func (p *page) Sleep(d time.Duration) Page {
	if p.err != nil {
		return p
	}
	p.log.Debugf("%s: waiting for %+v", p.batch, d)
	<-time.After(d)
	return p
}

func (p *page) Address() string {
	return p.page.URL()
}

func (p *page) checkError(method string) {
	if p.err == nil {
		return
	}

	tmpUUID := uuid.NewV4().String()

	if p.batch != "" {
		p.log.Errorf("batch '%s' error: %s failed: %+v [ %s ]", p.batch, method, p.err, tmpUUID)
	} else {
		p.log.Errorf("browser error: %s failed: %+v [ %s ]", method, p.err, tmpUUID)
	}

	if errors.Is(p.err, context.DeadlineExceeded) || errors.Is(p.err, context.Canceled) {
		return
	}

	if settings.StoreDumpOnError {
		if err := os.MkdirAll(p.dumpPath, os.ModePerm); err != nil {
			p.log.Warnf("cannot create dump directory: %s", err)
			return
		}
		p.Screenshot(path.Join(p.dumpPath, tmpUUID+".jpg")).TracePage(path.Join(p.dumpPath, tmpUUID+".html"))
	}
}

func (p *page) Close() {
	_ = p.page.Close()
}
