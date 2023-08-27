package deezer

import "github.com/RacoonMediaServer/rms-media-discovery/pkg/model"

type artistResponse struct {
	Name    string
	Picture string `json:"picture_medium"`
	Albums  uint   `json:"nb_album"`
}

type getArtistsResponse struct {
	Data []artistResponse
}

func (r getArtistsResponse) convert() []model.Music {
	var result []model.Music
	for _, a := range r.Data {
		artist := model.Artist{
			Name:       a.Name,
			PictureUrl: a.Picture,
			Albums:     a.Albums,
		}
		result = append(result, model.PackMusic(&artist))
	}
	return result
}
