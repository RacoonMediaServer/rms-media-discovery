package model

type Album struct {
	Title       string
	CoverUrl    string
	ReleaseDate string
	Genres      []string
	Tracks      uint
}

type AlbumResult struct {
	Album
	Artist string
}
