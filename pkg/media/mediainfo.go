package media

type VideoTrack struct {
	Width       int
	Height      int
	Codec       string
	AspectRatio string
}

type AudioTrack struct {
	Codec    string
	Language string
	Voice    string
}

type SubtitleTrack struct {
	Codec    string
	Language string
}

type Info struct {
	Format   string
	Video    []VideoTrack
	Audio    []AudioTrack
	Subtitle []SubtitleTrack
}
