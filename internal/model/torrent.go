package model

import "git.rms.local/RacoonMediaServer/rms-media-discovery/internal/media"

type Torrent struct {
	Title   string
	Link    string
	SizeMB  uint64
	Seeders uint
	Media   *media.Info
}
