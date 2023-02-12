package model

import (
	"context"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/pkg/heuristic"
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
	Info    heuristic.Info

	DetailLink string
	Downloader DownloadFunc
}

func (t Torrent) IsValid() bool {
	return t.Link != "" && t.Title != ""
}

type SearchQuery struct {
	Query  string
	Type   media.ContentType
	Limit  uint
	Year   *uint
	Season *uint
}
