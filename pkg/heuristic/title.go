package heuristic

import (
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/media"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type parseContext struct {
	title      string
	tokens     tokenList
	remove     []bool
	info       Info
	seasons    map[uint]struct{}
	subtitles  map[string]struct{}
	formatType media.ContentType
	codecType  media.ContentType
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
	ctx.parseSeasonCode()

	// парсим случай перечисления сезонов (сезоны 1-3)
	if !ctx.parseSeasons() {
		// парсим случаи, когда сезон указан отдельно (сезон 1, 1 сезон, season 1, 1-5 сезон etc)
		ctx.parseSeason()

		// случай, когда пишут 2 сезона, 5 сезонов
		ctx.parseShortSeason()

		// случай, когда просто указано количество серий вида [12 из 12]
		ctx.parseEpisodeRange()

		// случай, когда из названия можно угадать, что это сериал
		ctx.parseSeriesHint()
	}

	// вытаскиваем год (если диапазон, то только первую часть)
	ctx.parseYear()

	// 480p, 720p, FullHD...
	ctx.parseQuality()

	// для случая, если несколько фильмов в одной раздач
	ctx.parseTrilogy()

	// hdtvrip, dvdrip, webrip, ...
	ctx.parseRip()

	// mkv, avi, mp4...
	ctx.parseFormat()

	// h264, aac...
	ctx.parseCodec()

	// пробуем распознать языки субтитров
	ctx.parseSubtitles()

	// озвучку пробуем вытащить
	ctx.parseVoice()

	// исходя из распарсенного - пытаемся определить тип контента
	ctx.guessType()

	// удаляем лишние слова, чтобы определить вероятное название
	ctx.removeExtraWords()

	// пытаемся угадать название
	ctx.parseTitles()

	ctx.makeResult()

	return ctx.info
}

func mustParseUint(text string) uint {
	result, err := strconv.ParseUint(text, 10, 32)
	if err != nil {
		panic("must be digits" + err.Error())
	}
	return uint(result)
}

func (ctx *parseContext) parseSeasonCode() {
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

func (ctx *parseContext) parseSeasons() bool {
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

func (ctx *parseContext) parseSeason() {
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
			ctx.seasons[begin] = struct{}{}
			ctx.remove[pos] = true
			ctx.remove[found] = true

			if found > pos {
				end, ok := guessRangeEnd(ctx, found+1, begin)
				if ok {
					for i := begin + 1; i <= end; i++ {
						ctx.seasons[i] = struct{}{}
					}
				}
			} else if found > 0 && m.Match(ctx.tokens[found-1]) {
				end := begin
				begin = mustParseUint(ctx.tokens[found-1].Text)
				if begin < end {
					for i := begin; i < end; i++ {
						ctx.seasons[i] = struct{}{}
					}
					ctx.remove[found-1] = true
				}
			}
		}
	}
}

func (ctx *parseContext) parseShortSeason() {
	var m match
	m = &orMatch{
		Matches: []match{
			&wordMatch{Word: "сезона"},
			&wordMatch{Word: "сезонов"},
		},
	}
	pos := ctx.tokens.Find(m)
	if pos < 0 || pos == 0 {
		return
	}

	tok := ctx.tokens[pos-1]
	m = &regexMatch{Exp: seasonNumberRegex}
	if !m.Match(tok) {
		return
	}
	count := mustParseUint(tok.Text)
	for i := uint(1); i <= count; i++ {
		ctx.seasons[i] = struct{}{}
	}
	ctx.remove[pos] = true
	ctx.remove[pos-1] = true
}

func (ctx *parseContext) parseEpisodeRange() {
	if len(ctx.seasons) != 0 {
		return
	}

	separators := ctx.tokens.FindAll(&wordMatch{Word: "из"})
	for _, i := range separators {
		if i > 0 && i < len(ctx.tokens)-1 {
			prev := ctx.tokens[i-1]
			next := ctx.tokens[i+1]
			if prev.IsDigital() && next.IsDigital() {
				ctx.seasons[1] = struct{}{}
				ctx.remove[i-1] = true
				ctx.remove[i] = true
				ctx.remove[i+1] = true
				return
			}
		}
	}
}

func (ctx *parseContext) parseSeriesHint() {
	if len(ctx.seasons) != 0 {
		return
	}

	hints := orMatch{Matches: []match{
		wordMatch{Word: "ova"},
		wordMatch{Word: "тв"},
	}}

	if ctx.tokens.Find(&hints) >= 0 {
		ctx.seasons[1] = struct{}{}
	}
}

func (ctx *parseContext) parseYear() {
	years := ctx.tokens.FindAll(&regexMatch{Exp: yearRegex})
	for _, pos := range years {
		ctx.remove[pos] = true
	}
	if len(years) >= 1 {
		ctx.info.Year = mustParseUint(ctx.tokens[years[0]].Text)
	}
}

func (ctx *parseContext) makeResult() {
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

func (ctx *parseContext) parseQuality() {
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

func (ctx *parseContext) parseTrilogy() {
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

func (ctx *parseContext) parseRip() {
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

func (ctx *parseContext) parseFormat() {
	formats := map[string]media.ContentType{
		"mkv":  media.Movies,
		"mp4":  media.Movies,
		"webm": media.Movies,
		"avi":  media.Movies,
		"mpg":  media.Movies,
		"wmv":  media.Movies,
		"ogm":  media.Movies,
		"flv":  media.Movies,
		"mp3":  media.Music,
		"m4a":  media.Music,
		"flac": media.Music,
		"alac": media.Music,
		"ogg":  media.Music,
		"mka":  media.Music,
		"opus": media.Music,
	}
	m := &orMatch{}

	for format, t := range formats {
		if t == media.Movies {
			m.Matches = append(m.Matches, &wordMatch{format})
		}
	}
	for format, t := range formats {
		if t == media.Music {
			m.Matches = append(m.Matches, &wordMatch{format})
		}
	}

	pos := ctx.tokens.Find(m)
	if pos < 0 {
		return
	}
	ctx.remove[pos] = true
	ctx.info.Format = ctx.tokens[pos].Text
	ctx.formatType = formats[ctx.info.Format]
}

func (ctx *parseContext) parseCodec() {
	codecs := map[string]media.ContentType{
		"x264":   media.Movies,
		"h264":   media.Movies,
		"264":    media.Movies,
		"265":    media.Movies,
		"h265":   media.Movies,
		"hevc":   media.Movies,
		"avc":    media.Movies,
		"av1":    media.Movies,
		"vc1":    media.Movies,
		"vp8":    media.Movies,
		"vp9":    media.Movies,
		"mpeg":   media.Movies,
		"mpegts": media.Movies,
		"mp3":    media.Music,
		"flac":   media.Music,
		"alac":   media.Music,
		"aac":    media.Music,
		"ac3":    media.Music,
		"ogg":    media.Music,
	}
	m := &orMatch{}

	for format, t := range codecs {
		if t == media.Movies {
			m.Matches = append(m.Matches, &wordMatch{format})
		}
	}
	for format, t := range codecs {
		if t == media.Music {
			m.Matches = append(m.Matches, &wordMatch{format})
		}
	}

	pos := ctx.tokens.Find(m)
	if pos < 0 {
		return
	}
	ctx.remove[pos] = true
	ctx.info.Codec = ctx.tokens[pos].Text
	ctx.codecType = codecs[ctx.info.Format]
}

func (ctx *parseContext) guessType() {
	if len(ctx.seasons) != 0 || ctx.info.Rip != "" || ctx.info.Quality != media.QualityUnd || ctx.formatType == media.Movies || ctx.codecType == media.Movies {
		ctx.info.Type = media.Movies
		return
	}

	musicHints := &orMatch{
		Matches: []match{
			&wordMatch{"дискография"},
			&wordMatch{"discography"},
			&wordMatch{"ost"},
		},
	}
	if ctx.tokens.Find(musicHints) >= 0 {
		ctx.info.Type = media.Music
		return
	}

	if ctx.formatType == media.Music || ctx.codecType == media.Music {
		ctx.info.Type = media.Music
		return
	}
}

func (ctx *parseContext) removeExtraWords() {
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

func (ctx *parseContext) guessTitleLength() int {
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

func (ctx *parseContext) parseTitles() {
	l := ctx.guessTitleLength()
	tokens := crop(ctx, l)

	var result tokenList
	for _, t := range tokens {
		if t.SeqStart && len(result) != 0 {
			ctx.info.Titles = append(ctx.info.Titles, result.String())
			result = nil
		}
		if t.Text == "трилогия" || t.Text == "trilogy" || t.Text == "ova" || t.Text == "тв" {
			continue
		}
		result.Push(t)
	}
	if len(result) != 0 {
		ctx.info.Titles = append(ctx.info.Titles, result.String())
	}
}

func (ctx *parseContext) parseSubtitles() {
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

func (ctx *parseContext) parseVoice() {
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
