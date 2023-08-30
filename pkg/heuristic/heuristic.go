package heuristic

import (
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/media"
)

// Info contains all recognized information about the torrent
type Info struct {
	// Titles is a list of possible titles of content
	Titles []string

	// Seasons is a list of possible season numbers, which the torrent contains
	Seasons []uint

	// Year of movie creation (zero if don't recognized)
	Year uint

	// Quality is possible quality of video (resolution)
	Quality media.Quality

	// Trilogy means the torrent contains of a few movies
	Trilogy bool

	// Rip
	Rip string

	// Type represents possible type of torrent's content
	Type media.ContentType

	// Format means media container
	Format string

	// Codec
	Codec string

	// Voice contains info about voice acting (not parsed totally at all, just caught)
	Voice string

	// Subtitles is a list of subtitles language codes
	Subtitles []string
}
