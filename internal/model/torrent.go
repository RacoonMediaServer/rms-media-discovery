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

type SearchTypeHint int

const (
	SearchType_Movies SearchTypeHint = iota
	SearchType_Music
	SearchType_Books
	SearchType_Other
)

type SearchQuery struct {
	Query    string
	Hint     SearchTypeHint
	Limit    uint
	Detailed bool
}
