package torrents

import (
	"context"
	"errors"
)

func (s *Service) Download(ctx context.Context, link string) ([]byte, error) {
	val, ok := s.links.Load(link)
	if !ok {
		return nil, ErrExpiredDownloadLink
	}

	dl, ok := val.(*downloadLink)
	if !ok {
		return nil, errors.New("cannot resolve link properly")
	}

	return dl.downloader(ctx)
}
