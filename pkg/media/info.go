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

func (m *Info) VideoCount() int {
	return len(m.Video)
}

func (m *Info) CreateVideoTrack() int {
	m.Video = append(m.Video, VideoTrack{})
	return len(m.Video) - 1
}

func (m *Info) CreateAudioTrack() int {
	m.Audio = append(m.Audio, AudioTrack{})
	return len(m.Audio) - 1
}

func (m *Info) CreateSubtitleTrack() int {
	m.Subtitle = append(m.Subtitle, SubtitleTrack{})
	return len(m.Subtitle) - 1
}
