package model

type MovieType string

const (
	MovieType_TvSeries MovieType = "tv-series"
	MovieType_Movie    MovieType = "film"
)

type Movie struct {
	ID            string
	Description   string
	Genres        []string
	Poster        string
	Preview       string
	Rating        float32
	Seasons       uint
	Title         string
	OriginalTitle string
	Type          MovieType
	Year          uint
}
