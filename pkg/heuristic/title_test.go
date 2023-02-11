package heuristic

import (
	"git.rms.local/RacoonMediaServer/rms-media-discovery/pkg/media"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/pkg/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseTitle(t *testing.T) {
	type testCase struct {
		Title    string
		Expected Info
	}
	testCases := []testCase{
		{
			Title: "Паранормальный Веллингтон / Paranormal Unit / Wellington Paranormal [S01-03] (2018-2021) BDRip-HEVC 1080p | Кубик в кубе",
			Expected: Info{
				Titles:  []string{"Паранормальный Веллингтон", "Paranormal Unit", "Wellington Paranormal"},
				Seasons: []uint{1, 2, 3},
				Year:    2018,
				Quality: media.Quality1080p,
				Trilogy: false,
				Rip:     "bdrip",
				Type:    model.Movies,
				Format:  "",
				Voice:   "Кубик в кубе",
			},
		},
		{
			Title: "Паранормальный Веллингтон / Paranormal Unit / Wellington Paranormal [S02] (2019) WEBRip | Good People",
			Expected: Info{
				Titles:  []string{"Паранормальный Веллингтон", "Paranormal Unit", "Wellington Paranormal"},
				Seasons: []uint{2},
				Year:    2019,
				Quality: 0,
				Trilogy: false,
				Rip:     "webrip",
				Type:    model.Movies,
				Format:  "",
				Voice:   "Good People",
			},
		},
		{
			Title: "Чем мы заняты в тени / What We Do in the Shadows / Сезон: 1-4 / Серии: 1-40 из 40 (Джемейн Клемент, Джейсон Уолинер, Тайка Вайтити) [2019-2022, США, Комедия, ужасы, WEB-DLRip] MVO (LostFilm) + Original + (Rus, Eng)",
			Expected: Info{
				Titles:  []string{"Чем мы Заняты в Тени", "What we do in The Shadows"},
				Seasons: []uint{1, 2, 3, 4},
				Year:    2019,
				Quality: 0,
				Trilogy: false,
				Rip:     "webrip",
				Type:    model.Movies,
				Format:  "",
				Voice:   "Mvo Lostfilm Original Rus Eng",
			},
		},
		{
			Title: "Чем мы заняты в тени / What We Do in the Shadows / Сезон: 4 / Серии: 1-10 из 10 (Яна Горская, Кайл Ньюачек, Тиг Фонг, Д.Дж. Стипсен) [2022, США, Комедия, ужасы, WEB-DL 1080p] 3 x MVO (HDRezka Studio | Ozz | LostFilm) + Ukr (Ozz) + Original + Sub (Rus, Ukr, Eng)",
			Expected: Info{
				Titles:    []string{"Чем мы Заняты в Тени", "What we do in The Shadows"},
				Seasons:   []uint{4},
				Year:      2022,
				Quality:   media.Quality1080p,
				Trilogy:   false,
				Rip:       "webrip",
				Type:      model.Movies,
				Format:    "",
				Voice:     "Mvo Hdrezka Studio Ozz Lostfilm Ukr Ozz Original",
				Subtitles: []string{"en", "ru"},
			},
		},
		{
			Title: "The.Big.Bang.Theory.Season.11.Complete.720p.WEB-DL.x264.AAC",
			Expected: Info{
				Titles:  []string{"The Big Bang Theory"},
				Seasons: []uint{11},
				Year:    0,
				Quality: media.Quality720p,
				Trilogy: false,
				Rip:     "webrip",
				Type:    model.Movies,
				Format:  "",
			},
		},
		{
			Title: "The Big Bang Theory Season 1 2 3 4 5 6 7 8 9 10 - threesixtyp",
			Expected: Info{
				Titles:  []string{"The Big Bang Theory"},
				Seasons: []uint{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
				Year:    0,
				Quality: 0,
				Trilogy: false,
				Rip:     "",
				Type:    model.Movies,
				Format:  "",
			},
		},
		{
			Title: "Матрица: Трилогия / The Matrix: Trilogy [1999-2003, AC3, NTSC] [Open Matte] MVO + DVO (Гланц Королёва, Turkf) + AVO (CDV, Гаврилов, Живов, Кашкин, Пучков ) + VO (Есарев, 1+1) + Sub (Rus, Eng, Ukr) + Original Eng",
			Expected: Info{
				Titles:    []string{"Матрица", "The Matrix"},
				Seasons:   nil,
				Year:      1999,
				Quality:   0,
				Trilogy:   true,
				Rip:       "",
				Type:      model.Other,
				Format:    "",
				Voice:     "Mvo Dvo Гланц Королёва Turkf Avo Cdv Гаврилов Живов Кашкин Пучков vo Есарев 1 1",
				Subtitles: []string{"en", "ru"},
			},
		},
	}

	for i, c := range testCases {
		result := ParseTitle(c.Title)
		assert.Equal(t, c.Expected, result, "t#%d", i)
	}
}
