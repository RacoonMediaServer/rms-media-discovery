package model

import (
	"context"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/pkg/media"
)

// DownloadFunc is a function which can download the torrent
type DownloadFunc func(ctx context.Context) ([]byte, error)

// Torrent is an internal representation of torrent record
type Torrent struct {
	Title   string
	Link    string
	SizeMB  uint64
	Seeders uint
	Media   *media.Info

	DetailLink string
	Downloader DownloadFunc
}

type ContentType int

const (
	Movies ContentType = iota
	Music
	Books
	Other
)

// OrderByFunc is a func for sort results
type OrderByFunc func(a, b *Torrent) bool

type SearchQuery struct {
	Query    string
	Type     ContentType
	Limit    uint
	Detailed bool
	Year     *uint
	Season   *uint
	OrderBy  OrderByFunc
}

func OrderBySeeders(a, b *Torrent) bool {
	return a.Seeders > b.Seeders
}

func OrderByTitle(a, b *Torrent) bool {
	return a.Title < b.Title
}

func OrderBySize(a, b *Torrent) bool {
	return a.SizeMB > b.SizeMB
}
