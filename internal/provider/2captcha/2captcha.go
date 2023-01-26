package _captcha

import (
	"context"
	"encoding/base64"
	"fmt"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/model"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/provider"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/requester"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/utils"
	api2captcha "github.com/2captcha/2captcha-go"
	"github.com/apex/log"
)

type captchaSolver struct {
	log    *log.Entry
	access model.AccessProvider
	r      requester.Requester
}

func (c captchaSolver) ID() string {
	return "2captcha"
}

func (c captchaSolver) Solve(ctx context.Context, captcha provider.Captcha) (string, error) {
	l := utils.LogFromContext(ctx, "2captcha", c.log)
	l.Info("Captcha resolving requested")
	content, err := c.r.Download(ctx, captcha.Url)
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

func NewSolver(access model.AccessProvider) provider.CaptchaSolver {
	s := &captchaSolver{access: access, log: log.WithField("from", "2captcha")}
	s.r = requester.New(s)
	return s
}
