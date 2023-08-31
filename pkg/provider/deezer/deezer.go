package deezer

import (
	"context"
	"fmt"
	"github.com/RacoonMediaServer/rms-media-discovery/internal/utils"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/model"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/pipeline"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/provider"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/requester"
	"github.com/antzucaro/matchr"
	"net/url"
	"strings"
	"unicode/utf8"
)

type deezerProvider struct {
	p pipeline.Pipeline
	r requester.Requester
}

const deezerEndpoint = "https://api.deezer.com"

func (d deezerProvider) ID() string {
	return "deezer"
}

func (d deezerProvider) SearchMusic(ctx context.Context, query string, limit uint, searchType model.MusicSearchType) ([]model.Music, error) {
	l := utils.LogFromContext(ctx, "deezer")
	var result []model.Music
	var err error

	if searchType == model.SearchAny || searchType == model.SearchArtist {
		artists := getArtistsResponse{}
		err = d.r.Get(ctx, fmt.Sprintf("%s/search/artist?q=%s&limit=%d", deezerEndpoint, url.PathEscape(query), limit), &artists)
		if err != nil {
			l.Errorf("Search artists failed: %s", err)
			return nil, err
		}
		result = append(result, artists.convert()...)
	}

	if searchType == model.SearchAny || searchType == model.SearchAlbum {
		albums := &getAlbumsResponse{}
		err = d.r.Get(ctx, fmt.Sprintf("%s/search/album?q=%s&limit=%d", deezerEndpoint, url.PathEscape(query), limit), &albums)
		if err != nil {
			l.Errorf("Search albums failed: %s", err)
			return nil, err
		}
		result = append(result, albums.convert()...)
	}

	if searchType == model.SearchAny || searchType == model.SearchTrack {
		tracks := &searchTrackResponse{}
		err = d.r.Get(ctx, fmt.Sprintf("%s/search/track?q=%s&limit=%d", deezerEndpoint, url.PathEscape(query), limit), &tracks)
		if err != nil {
			l.Errorf("Search tracks failed: %s", err)
			return nil, err
		}
		result = append(result, tracks.convert()...)
	}

	target := strings.ToLower(query)
	filtered := make([]model.Music, 0, len(result))
	for _, res := range result {
		distance := matchr.Levenshtein(strings.ToLower(res.Title()), target)
		if distance < int(float32(utf8.RuneCountInString(target))*0.75) {
			filtered = append(filtered, res)
		}
	}
	if len(filtered) == 0 {
		filtered = result
	}

	return filtered, nil
}

func NewProvider() provider.MusicInfoProvider {
	p := &deezerProvider{
		p: pipeline.Open(pipeline.Settings{Id: "deezer"}),
	}
	p.r = requester.New(p)
	return p
}
