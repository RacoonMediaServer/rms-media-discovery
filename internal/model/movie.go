package model

type Movie struct {
	ID          string
	Description string
	Genres      []string
	Poster      string
	Preview     string
	Rating      float32
	Seasons     uint
	Title       string
	Type        string
	Year        uint
}
