package tmdb

import (
	"context"
	"errors"
	"github.com/ryanbradynd05/go-tmdb"
	"strings"
)

type requestFunc func(api *tmdb.TMDb) (interface{}, error)

func (p *tmdbProvider) request(ctx context.Context, f requestFunc) (interface{}, error) {
	for {
		resp, err := p.p.Do(ctx, func() (interface{}, error) {
			token, err := p.access.GetApiKey(p.ID())
			if err != nil {
				return nil, err
			}
			config := tmdb.Config{APIKey: token.Key}
			api := tmdb.Init(config)

			result, err := f(api)
			if err != nil {
				if strings.Index(err.Error(), "Invalid API key") >= 0 {
					p.access.MarkUnaccesible(token.AccountId)
					return nil, errBadAccount
				}
				return nil, err
			}
			return result, nil
		})

		if err != nil {
			if errors.Is(err, errBadAccount) {
				continue
			}

			return nil, err
		}

		return resp, err
	}
}
