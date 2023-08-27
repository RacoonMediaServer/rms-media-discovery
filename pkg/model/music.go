package model

type Music struct {
	content interface{}
}

func (m Music) IsArtist() bool {
	_, ok := m.content.(*Artist)
	return ok
}

func (m Music) IsAlbum() bool {
	_, ok := m.content.(*AlbumResult)
	return ok
}

func (m Music) IsTrack() bool {
	_, ok := m.content.(*TrackResult)
	return ok
}

func (m Music) AsArtist() *Artist {
	return m.content.(*Artist)
}

func (m Music) AsAlbum() *AlbumResult {
	return m.content.(*AlbumResult)
}

func (m Music) AsTrack() *TrackResult {
	return m.content.(*TrackResult)
}

func (m Music) Title() string {
	switch item := m.content.(type) {
	case *Artist:
		return item.Name
	case *AlbumResult:
		return item.Title
	case *TrackResult:
		return item.Title
	default:
		return ""
	}
}
func PackMusic[T Artist | AlbumResult | TrackResult](item *T) Music {
	return Music{content: item}
}
