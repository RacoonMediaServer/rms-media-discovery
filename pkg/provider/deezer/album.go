package deezer

import "github.com/RacoonMediaServer/rms-media-discovery/pkg/model"

type albumResponse struct {
	Title  string
	Cover  string `json:"cover_medium"`
	Tracks uint   `json:"nb_tracks"`
	Artist artistResponse
}

type getAlbumsResponse struct {
	Data []albumResponse
}

func (r getAlbumsResponse) convert() []model.Music {
	var result []model.Music
	for _, a := range r.Data {
		album := model.AlbumResult{
			Album: model.Album{
				Title:       a.Title,
				CoverUrl:    a.Cover,
				ReleaseDate: "",
				Genres:      nil,
				Tracks:      a.Tracks,
			},
			Artist: a.Artist.Name,
		}
		result = append(result, model.PackMusic(&album))
	}
	return result
}
