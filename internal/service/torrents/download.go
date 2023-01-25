package torrents

import (
	"context"
	"errors"
)

func (s *service) Download(ctx context.Context, link string) ([]byte, error) {
	s.cleanExpiredLinks()

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
