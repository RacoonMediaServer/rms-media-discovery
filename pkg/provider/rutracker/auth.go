package rutracker

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/model"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/navigator"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/provider"
	"github.com/apex/log"
)

func authorize(ctx context.Context, cred model.Credentials, solver provider.CaptchaSolver) (cookies []*http.Cookie, err error) {
	var nav navigator.Navigator
	var p navigator.Page
	l := log.WithField("from", "rutracker").WithField("account", cred.AccountId)

	nav, err = navigator.New("rutracker")
	if err != nil {
		err = fmt.Errorf("cannot create session browser: %w", err)
		return
	}
	defer nav.Close()

	p, err = nav.NewPage(l, ctx)
	if err != nil {
		err = fmt.Errorf("create browser page failed: %w", err)
		return
	}
	defer p.Close()

	err = p.Batch("authorization").
		Goto("https://rutracker.org/forum/login.php?redirect=search.php").
		FetchContent().
		Error()

	if err != nil {
		return
	}

	var captcha struct {
		required bool
		url      string
		code     string
	}

	doc := p.Document()
	captcha.url, captcha.required = doc.Find("#login-form-full > table > tbody > tr:nth-child(2) > td > div > table > tbody > tr:nth-child(3) > td:nth-child(2) > div:nth-child(1) > img").
		First().
		Attr("src")

	if captcha.required {
		l.Warn("Captcha is required")

		code := ""
		code, err = solver.Solve(ctx, provider.Captcha{
			Url:           captcha.url,
			CaseSensitive: false,
			MinLength:     4,
			MaxLength:     6,
		})
		if err != nil {
			err = fmt.Errorf("cannot solve captcha: %w", err)
			return
		}

		err = p.Batch("authorization with captcha").
			Goto("https://rutracker.org/forum/login.php?redirect=search.php").
			Type(`#login-form-full > table > tbody > tr:nth-child(2) > td > div > table > tbody > tr:nth-child(1) > td:nth-child(2) > input[type=text]`, cred.Login).
			Type(`#login-form-full > table > tbody > tr:nth-child(2) > td > div > table > tbody > tr:nth-child(2) > td:nth-child(2) > input[type=password]`, cred.Password).
			Type(`#login-form-full > table > tbody > tr:nth-child(2) > td > div > table > tbody > tr:nth-child(3) > td:nth-child(2) > div:nth-child(2) > input.reg-input`, code).
			Click(`#login-form-full > table > tbody > tr:nth-child(2) > td > div > table > tbody > tr:nth-child(4) > td > input`).
			FetchContent().
			Error()

	} else {
		err = p.Batch("authorization").
			Type(`#login-form-full > table > tbody > tr:nth-child(2) > td > div > table > tbody > tr:nth-child(1) > td:nth-child(2) > input[type=text]`, cred.Login).
			Type(`#login-form-full > table > tbody > tr:nth-child(2) > td > div > table > tbody > tr:nth-child(2) > td:nth-child(2) > input[type=password]`, cred.Password).
			Click(`#login-form-full > table > tbody > tr:nth-child(2) > td > div > table > tbody > tr:nth-child(4) > td > input`).
			FetchContent().
			Error()
	}

	if err != nil {
		return
	}

	doc = p.Document()

	authorized := false
	doc.Find("#logged-in-username").Each(func(i int, selection *goquery.Selection) {
		authorized = true
	})

	if !authorized {
		p.RaiseError(errors.New("cannot parse account logged"))
		err = errBadAccount
	}

	cookies, err = nav.GetCookies()
	return
}
