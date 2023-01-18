package model

type MovieType string

const (
	MovieType_TvSeries MovieType = "tv-series"
	MovieType_Movie    MovieType = "movie"
)

type Movie struct {
	ID          string
	Description string
	Genres      []string
	Poster      string
	Preview     string
	Rating      float32
	Seasons     uint
	Title       string
	Type        MovieType
	Year        uint
}
