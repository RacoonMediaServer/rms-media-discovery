package deezer

import "github.com/RacoonMediaServer/rms-media-discovery/pkg/model"

type trackResponse struct {
	Title  string
	Album  albumResponse
	Artist artistResponse
}

type searchTrackResponse struct {
	Data []trackResponse
}

func (r searchTrackResponse) convert() []model.Music {
	var result []model.Music
	for _, t := range r.Data {
		song := model.TrackResult{
			Track: model.Track{
				Title:    t.Title,
				Position: 0,
			},
			Artist:   t.Artist.Name,
			Album:    t.Album.Title,
			CoverUrl: t.Album.Cover,
		}
		result = append(result, model.PackMusic(&song))
	}
	return result
}
