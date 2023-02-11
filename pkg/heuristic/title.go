package heuristic

import (
	"git.rms.local/RacoonMediaServer/rms-media-discovery/pkg/media"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type parseContext struct {
	title     string
	tokens    tokenList
	remove    []bool
	info      Info
	seasons   map[uint]struct{}
	subtitles map[string]struct{}
}

var (
	seasonCodeRegex   = regexp.MustCompile(`^s(\d\d)$`)
	seasonNumberRegex = regexp.MustCompile(`^(\d\d?)$`)
	yearRegex         = regexp.MustCompile(`^\d\d\d\d$`)
	ripRegex          = regexp.MustCompile(`^\w+rip$`)
)

// ParseTitle parses torrent's title and extract some useful info, such as season number, year, quality etc
func ParseTitle(title string) Info {
	// разбиваем текст на токены
	tokens := parse(title)

	ctx := parseContext{
		title:     title,
		tokens:    tokens,
		remove:    make([]bool, len(tokens)),
		seasons:   make(map[uint]struct{}),
		subtitles: make(map[string]struct{}),
	}

	// парсим случаи, когда в тексте есть 'S01' или 'S01' 'S02' или 'S01-12'
	parseSeasonCode(&ctx)

	// парсим случай перечисления сезонов (сезоны 1-3)
	if !parseSeasons(&ctx) {
		// парсим случаи, когда сезон указан отдельно (сезон 1, 1 сезон, season 1, etc)
		parseSeason(&ctx)
	}

	// вытаскиваем год (если диапазон, то только первую часть)
	parseYear(&ctx)

	// 480p, 720p, FullHD...
	parseQuality(&ctx)

	// для случая, если несколько фильмов в одной раздач
	parseTrilogy(&ctx)

	// hdtvrip, dvdrip, webrip, ...
	parseRip(&ctx)

	// mkv, avi, mp4...
	parseFormat(&ctx)

	// пробуем распознать языки субтитров
	parseSubtitles(&ctx)

	// озвучку пробуем вытащить
	parseVoice(&ctx)

	// исходя из распарсенного - пытаемся определить тип контента
	guessType(&ctx)

	// удаляем лишние слова, чтобы определить вероятное название
	removeExtraWords(&ctx)

	// пытаемся угадать название
	parseTitles(&ctx)

	makeResult(&ctx)

	return ctx.info
}

func mustParseUint(text string) uint {
	result, err := strconv.ParseUint(text, 10, 32)
	if err != nil {
		panic("must be digits" + err.Error())
	}
	return uint(result)
}

func parseSeasonCode(ctx *parseContext) {
	m := regexMatch{Exp: seasonCodeRegex}
	pos := ctx.tokens.Find(m)
	if pos > -1 {
		tok := &ctx.tokens[pos]
		matches := m.Exp.FindStringSubmatch(tok.Text)
		if matches != nil {
			begin := mustParseUint(matches[1])
			// тут может быть диапазон
			if pos < len(ctx.tokens)-1 {
				tok = &ctx.tokens[pos+1]
				if tok.IsDigital() && len(tok.Text) < 3 {
					end := mustParseUint(ctx.tokens[pos+1].Text)
					if begin < end {
						for i := begin; i <= end; i++ {
							ctx.seasons[i] = struct{}{}
						}
						ctx.remove[pos+1] = true
					}
				}
			}
			ctx.seasons[begin] = struct{}{}
			ctx.remove[pos] = true
		}
	}
}

func parseSeasons(ctx *parseContext) bool {
	splitSeasonMatch := &orMatch{
		Matches: []match{
			&wordMatch{Word: "сезоны"},
			&wordMatch{Word: "seasons"},
		},
	}
	pos := ctx.tokens.Find(splitSeasonMatch)
	if pos < 0 || pos >= len(ctx.tokens)-1 {
		return false
	}
	ctx.remove[pos] = true

	m := &regexMatch{Exp: seasonNumberRegex}
	if !m.Match(ctx.tokens[pos+1]) {
		return false
	}
	ctx.remove[pos+1] = true

	begin := mustParseUint(ctx.tokens[pos+1].Text)
	end, ok := guessRangeEnd(ctx, pos+2, begin)
	if !ok {
		ctx.seasons[begin] = struct{}{}
		return true
	}

	for i := begin + 1; i <= end; i++ {
		ctx.seasons[i] = struct{}{}
	}
	return true
}

func guessRangeEnd(ctx *parseContext, pos int, begin uint) (uint, bool) {
	if pos >= len(ctx.tokens) {
		return 0, false
	}
	m := &regexMatch{Exp: seasonNumberRegex}
	end := uint(0)
	for i := pos; i < len(ctx.tokens); i++ {
		if !m.Match(ctx.tokens[i]) {
			break
		}
		end = mustParseUint(ctx.tokens[i].Text)
		ctx.remove[i] = true
	}
	if end == 0 || begin > end {
		return 0, false
	}
	return end, true
}

func parseSeason(ctx *parseContext) {
	splitSeasonMatch := &orMatch{
		Matches: []match{
			&wordMatch{Word: "сезон"},
			&wordMatch{Word: "season"},
			&wordMatch{Word: "sezon"},
		},
	}
	pos := ctx.tokens.Find(splitSeasonMatch)
	if pos > -1 {
		m := regexMatch{Exp: regexp.MustCompile(`^\d\d?$`)}
		found := -1
		if pos < len(ctx.tokens)-1 && m.Match(ctx.tokens[pos+1]) {
			found = pos + 1
		} else if pos > 0 && m.Match(ctx.tokens[pos-1]) {
			found = pos - 1
		}
		if found < 0 && pos > 1 && ctx.tokens[pos-1].Text == "й" && m.Match(ctx.tokens[pos-2]) {
			found = pos - 2
		}

		if found > -1 {
			begin := mustParseUint(ctx.tokens[found].Text)
			ctx.info.Seasons = append(ctx.info.Seasons, begin)
			ctx.remove[pos] = true
			ctx.remove[found] = true

			end, ok := guessRangeEnd(ctx, found+1, begin)
			if ok {
				for i := begin + 1; i <= end; i++ {
					ctx.seasons[i] = struct{}{}
				}
			}
		}
	}
}

func parseYear(ctx *parseContext) {
	years := ctx.tokens.FindAll(&regexMatch{Exp: yearRegex})
	for _, pos := range years {
		ctx.remove[pos] = true
	}
	if len(years) >= 1 {
		ctx.info.Year = mustParseUint(ctx.tokens[years[0]].Text)
	}
}

func makeResult(ctx *parseContext) {
	for k, _ := range ctx.seasons {
		ctx.info.Seasons = append(ctx.info.Seasons, k)
	}
	sort.Slice(ctx.info.Seasons, func(i, j int) bool {
		return ctx.info.Seasons[i] < ctx.info.Seasons[j]
	})
	for k, _ := range ctx.subtitles {
		ctx.info.Subtitles = append(ctx.info.Subtitles, k)
	}
	sort.Slice(ctx.info.Subtitles, func(i, j int) bool {
		return ctx.info.Subtitles[i] < ctx.info.Subtitles[j]
	})
	ctx.info.Voice = strings.TrimSpace(ctx.info.Voice)
}

func parseQuality(ctx *parseContext) {
	qmap := map[string]media.Quality{
		"480p":   media.Quality480p,
		"720p":   media.Quality720p,
		"1080p":  media.Quality1080p,
		"2160p":  media.Quality2160p,
		"4k":     media.Quality2160p,
		"uhd":    media.Quality2160p,
		"fullhd": media.Quality1080p,
		"hd":     media.Quality720p,
		"hdtv":   media.Quality720p,
	}

	m := &orMatch{}
	for k, _ := range qmap {
		m.Matches = append(m.Matches, &wordMatch{Word: k})
	}

	pos := ctx.tokens.Find(m)
	if pos < 0 {
		return
	}

	ctx.info.Quality = qmap[ctx.tokens[pos].Text]
	ctx.remove[pos] = true
}

func parseTrilogy(ctx *parseContext) {
	m := orMatch{
		Matches: []match{
			&wordMatch{"трилогия"},
			&wordMatch{"trilogy"},
		},
	}
	pos := ctx.tokens.FindAll(m)
	if len(pos) == 0 {
		return
	}

	ctx.info.Trilogy = true
}

func parseRip(ctx *parseContext) {
	m := &orMatch{
		Matches: []match{
			&regexMatch{Exp: ripRegex},
			&wordMatch{Word: "web"},
			&wordMatch{Word: "dl"},
		},
	}

	pos := ctx.tokens.Find(m)
	if pos < 0 {
		return
	}
	ctx.remove[pos] = true
	ctx.info.Rip = ctx.tokens[pos].Text
	switch ctx.info.Rip {
	case "dlrip":
		fallthrough
	case "dl":
		fallthrough
	case "web":
		ctx.info.Rip = "webrip"
	}
}

func parseFormat(ctx *parseContext) {
	m := &orMatch{
		Matches: []match{
			&wordMatch{"mkv"},
			&wordMatch{"mp4"},
			&wordMatch{"avi"},
		},
	}

	pos := ctx.tokens.Find(m)
	if pos < 0 {
		return
	}
	ctx.remove[pos] = true
	ctx.info.Format = ctx.tokens[pos].Text
}

func guessType(ctx *parseContext) {
	if len(ctx.seasons) != 0 || ctx.info.Rip != "" || ctx.info.Quality != media.QualityUnd || ctx.info.Format != "" {
		ctx.info.Type = media.Movies
		return
	}
}

func removeExtraWords(ctx *parseContext) {
	matched := ctx.tokens.FindAll(
		&orMatch{
			Matches: []match{
				&bracesMatch{},
				&regexMatch{Exp: regexp.MustCompile(`remux$`)},
			},
		},
	)

	for _, r := range matched {
		ctx.remove[r] = true
	}
}

func guessTitleLength(ctx *parseContext) int {
	for i, r := range ctx.remove {
		if r && i != 0 {
			if i == 1 && ctx.tokens[i-1].IsDigital() {
				continue
			}
			return i
		}
	}

	return len(ctx.remove)
}

func crop(ctx *parseContext, maxLength int) tokenList {
	result := make([]token, 0, maxLength)
	for i, t := range ctx.tokens {
		if !ctx.remove[i] {
			result = append(result, t)
		}
		if len(result) >= maxLength {
			return result
		}
	}

	return result
}

func parseTitles(ctx *parseContext) {
	l := guessTitleLength(ctx)
	tokens := crop(ctx, l)

	var result tokenList
	for _, t := range tokens {
		if t.SeqStart && len(result) != 0 {
			ctx.info.Titles = append(ctx.info.Titles, result.String())
			result = nil
		}
		if t.Text == "трилогия" || t.Text == "trilogy" {
			continue
		}
		result.Push(t)
	}
	if len(result) != 0 {
		ctx.info.Titles = append(ctx.info.Titles, result.String())
	}
}

func parseSubtitles(ctx *parseContext) {
	m := wordMatch{Word: "sub"}
	pos := ctx.tokens.Find(m)
	if pos < 0 {
		return
	}

	ctx.remove[pos] = true
	for i := pos + 1; i < len(ctx.tokens); i++ {
		content := ctx.tokens[i].Text
		code := ""
		switch content {
		case "rus":
			fallthrough
		case "ru":
			code = "ru"
		case "eng":
			fallthrough
		case "en":
			code = "en"
		}
		if code != "" {
			ctx.subtitles[code] = struct{}{}
		}
		ctx.remove[i] = true
	}
}

func parseVoice(ctx *parseContext) {
	m := &orMatch{
		Matches: []match{
			&wordMatch{"dub"},
			&wordMatch{"mvo"},
			&wordMatch{"dvo"},
			&wordMatch{"avo"},
			&wordMatch{"vo"},
		},
	}
	pos := ctx.tokens.Find(m)
	if pos > -1 {
		for i := pos; i < len(ctx.tokens); i++ {
			if !ctx.remove[i] {
				ctx.remove[i] = true
				ctx.info.Voice += ctx.tokens[i].String() + " "
			}
		}
		return
	}
	// на rutor.info озвучка записывается через вертикальную черту в конце строки
	pos = strings.LastIndex(ctx.title, "|")
	if pos > -1 {
		ctx.info.Voice = strings.TrimSpace(ctx.title[pos+1:])
	}
}
