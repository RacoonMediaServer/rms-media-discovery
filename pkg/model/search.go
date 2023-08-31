package model

type MusicSearchType int

const (
	SearchAny MusicSearchType = iota
	SearchArtist
	SearchAlbum
	SearchTrack
)
