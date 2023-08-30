package heuristic

import (
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/media"
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
				Type:    media.Movies,
				Format:  "",
				Codec:   "hevc",
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
				Type:    media.Movies,
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
				Type:    media.Movies,
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
				Type:      media.Movies,
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
				Type:    media.Movies,
				Codec:   "x264",
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
				Type:    media.Movies,
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
				Type:      media.Other,
				Format:    "",
				Codec:     "ac3",
				Voice:     "Mvo Dvo Гланц Королёва Turkf Avo Cdv Гаврилов Живов Кашкин Пучков vo Есарев 1 1",
				Subtitles: []string{"en", "ru"},
			},
		},
		{
			Title: "Теория Большого Взрыва / The Big Bang Theory [1-5 сезон] (Марк Сендровски, Джеймс Берроуз) [2007-2011, Комедия, HDTVRip] [MP4, 640x] (Кураж-Бамбей)",
			Expected: Info{
				Titles:    []string{"Теория Большого Взрыва", "The Big Bang Theory"},
				Seasons:   []uint{1, 2, 3, 4, 5},
				Year:      2007,
				Quality:   0,
				Trilogy:   false,
				Rip:       "hdtvrip",
				Type:      media.Movies,
				Format:    "mp4",
				Voice:     "",
				Subtitles: nil,
			},
		},
		{
			Title: "Гильдия / The Guild / Сезон: 2(полный)+Special s02 / (Джейн Селле Морган, Грег Бенсон) [2008-2009, США, Комедия, WEBRip] Озвучка, субтитры",
			Expected: Info{
				Titles:    []string{"Гильдия", "The Guild"},
				Seasons:   []uint{2},
				Year:      2008,
				Quality:   0,
				Trilogy:   false,
				Rip:       "webrip",
				Type:      media.Movies,
				Format:    "",
				Voice:     "",
				Subtitles: nil,
			},
		},
		{
			Title: "Гильдия / The Guild (Джейн Селле Морган, Грег Бенсон) [2 сезона] [22 серии] [2007, ситком, WEBRip] Sub",
			Expected: Info{
				Titles:    []string{"Гильдия", "The Guild"},
				Seasons:   []uint{1, 2},
				Year:      2007,
				Quality:   0,
				Trilogy:   false,
				Rip:       "webrip",
				Type:      media.Movies,
				Format:    "",
				Voice:     "",
				Subtitles: nil,
			},
		},
		{
			Title: "Ария - Ночь короче дня - 1995 - ремастер by Zvezdopad [FLAC (image+.cue), lossless]",
			Expected: Info{
				Titles:    []string{"Ария Ночь Короче Дня"},
				Year:      1995,
				Quality:   0,
				Type:      media.Music,
				Format:    "flac",
				Codec:     "flac",
				Voice:     "",
				Subtitles: nil,
			},
		},
		{
			Title: "Король и Шут - Дискография (18 релизов) - 1993-2013 (2018), AAC (tracks), 320 kbps (VBR)",
			Expected: Info{
				Titles:    []string{"Король и Шут Дискография"},
				Year:      1993,
				Quality:   0,
				Type:      media.Music,
				Format:    "",
				Codec:     "aac",
				Voice:     "",
				Subtitles: nil,
			},
		},
		{
			Title: "Iron Maiden - Senjutsu (2CD) (2021) Mp3 320kbps [PMEDIA] ⭐️",
			Expected: Info{
				Titles:    []string{"Iron Maiden Senjutsu"},
				Year:      2021,
				Quality:   0,
				Type:      media.Music,
				Format:    "mp3",
				Codec:     "mp3",
				Voice:     "",
				Subtitles: nil,
			},
		},
	}

	for i, c := range testCases {
		result := ParseTitle(c.Title)
		assert.Equal(t, c.Expected, result, "t#%d", i)
	}
}
