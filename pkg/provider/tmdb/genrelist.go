package tmdb

import (
	"context"
	"github.com/ryanbradynd05/go-tmdb"
	"strings"
	"sync"
)

type genreList struct {
	genres map[int]string
	mu     sync.RWMutex
}

func castGenreList[ID int32 | uint32 | uint](ids []ID) []int {
	result := make([]int, 0, len(ids))
	for i := range ids {
		result = append(result, int(ids[i]))
	}
	return result
}

func (p *tmdbProvider) initGenreList(ctx context.Context) {
	if p.g.isInitialized() {
		return
	}

	resp, err := p.request(ctx, func(api *tmdb.TMDb) (interface{}, error) {
		return api.GetMovieGenres(map[string]string{"language": "ru-RU"})
	})
	if err != nil {
		return
	}

	result := resp.(*tmdb.Genre)

	genres := make(map[int]string)
	for _, g := range result.Genres {
		genres[g.ID] = strings.ToLower(g.Name)
	}
	resp, err = p.request(ctx, func(api *tmdb.TMDb) (interface{}, error) {
		return api.GetTvGenres(map[string]string{"language": "ru-RU"})
	})
	if err != nil {
		return
	}
	result = resp.(*tmdb.Genre)
	for _, g := range result.Genres {
		genres[g.ID] = strings.ToLower(g.Name)
	}
	p.g.initialize(genres)
}

func (l *genreList) isInitialized() bool {
	l.mu.RLock()
	defer l.mu.RUnlock()

	return l.genres != nil
}

func (l *genreList) initialize(genres map[int]string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.genres = genres
}

func (l *genreList) get(ids []int) []string {
	l.mu.RLock()
	defer l.mu.RUnlock()

	if l.genres == nil {
		return nil
	}

	var genres []string
	for _, id := range ids {
		genre := l.genres[id]
		if genre != "" {
			genres = append(genres, genre)
		}
	}

	return genres
}
