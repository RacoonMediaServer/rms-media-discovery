package model

import (
	"context"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/media"
)

type DownloadFunc func(ctx context.Context) ([]byte, error)

type Torrent struct {
	Title   string
	Link    string
	SizeMB  uint64
	Seeders uint
	Media   *media.Info

	DetailLink string
	Downloader DownloadFunc
}
