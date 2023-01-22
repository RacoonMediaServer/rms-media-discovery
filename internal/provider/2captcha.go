package provider

import (
	"context"
	"encoding/base64"
	"fmt"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/model"
	api2captcha "github.com/2captcha/2captcha-go"
	"github.com/apex/log"
	"net/http"
)

type captchaSolver struct {
	log    *log.Entry
	access model.AccessProvider
}

func (c captchaSolver) Solve(ctx context.Context, captcha Captcha) (string, error) {
	content, err := download(c.log, http.Client{}, ctx, captcha.Url)
	if err != nil {
		return "", fmt.Errorf("download captcha failed: %w", err)
	}
	account, err := c.access.GetApiKey("2captcha")
	if err != nil {
		return "", err
	}

	cli := api2captcha.NewClient(account.Key)

	uCaptcha := api2captcha.Normal{
		Base64:        base64.StdEncoding.EncodeToString(content),
		CaseSensitive: captcha.CaseSensitive,
	}
	code, err := cli.Solve(uCaptcha.ToRequest())
	if err != nil {
		return "", err
	}

	return code, nil
}

func NewCaptchaSolver(access model.AccessProvider) CaptchaSolver {
	return &captchaSolver{access: access, log: log.WithField("from", "2captcha")}
}
