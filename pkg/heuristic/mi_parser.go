package heuristic

import (
	"bufio"
	"bytes"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/pkg/media"
	"golang.org/x/text/language"
	"golang.org/x/text/language/display"
	"regexp"
	"strconv"
	"strings"
)

var (
	formatRegex     = regexp.MustCompile(`^(Формат|Format)( видео)?: (.{1,20})`)
	videoRegex      = regexp.MustCompile(`^(Видео|Video)\s?#?(\d*): (.{1,100})`)
	audioRegex      = regexp.MustCompile(`^(Аудио|Audio)\s?#?(\d*): (.{1,100})`)
	resolutionRegex = regexp.MustCompile(`(\d+)x(\d+)`)
	subtitlesRegex  = regexp.MustCompile(`^Субтитры: (.{1,100})`)
)

// MediaInfoParser is heuristic parser, which can extract information about media from text
type MediaInfoParser struct {
	Info media.Info
}

// ParseMediaInfo parses amount of text and fills MediaInfoParser.Info
func (p *MediaInfoParser) ParseMediaInfo(text string) *media.Info {
	text = strings.ReplaceAll(text, string(rune(0xa0)), " ")
	text = strings.ReplaceAll(text, "<br />", "\n")
	// пробуем найти прям описание MediaInfo в тексте
	if p.parseMediaInfo(text) {
		return &p.Info
	}

	scanner := bufio.NewScanner(bytes.NewReader([]byte(text)))
	for scanner.Scan() {
		line := scanner.Text()
		if sm := formatRegex.FindStringSubmatch(line); sm != nil {
			p.Info.Format = sm[3]
		}
		p.parseVideo(line)
		p.parseAudio(line)
		p.parseSubtitles(line)
	}

	return &p.Info
}

func (p *MediaInfoParser) parseMediaInfo(text string) bool {
	if _, mediaInfo, ok := strings.Cut(text, "MediaInfo"); ok {
		p.Info = *parseMediaInfo(mediaInfo)
		return true
	}
	if _, mediaInfo, ok := strings.Cut(text, "MI\n"); ok {
		p.Info = *parseMediaInfo(mediaInfo)
		return true
	}

	return false
}

func (p *MediaInfoParser) parseVideo(line string) {
	if sm := videoRegex.FindStringSubmatch(line); sm != nil {
		videoLine := sm[3]
		id := p.Info.CreateVideoTrack()
		if sm = resolutionRegex.FindStringSubmatch(videoLine); sm != nil {
			w, _ := strconv.ParseInt(sm[1], 10, 32)
			h, _ := strconv.ParseInt(sm[2], 10, 32)
			p.Info.Video[id].Width, p.Info.Video[id].Height = int(w), int(h)
		}
	}
}

func parseLanguage(line string) (language.Tag, bool) {
	low := strings.ToLower(line)
	isMatched := func(n display.Namer, tag language.Tag) bool {
		try := strings.ToLower(n.Name(tag))
		return strings.Index(low, try) != -1
	}

	en := display.English.Tags()
	ru := display.Russian.Tags()

	for _, l := range Languages {
		if isMatched(en, l) || isMatched(ru, l) || isMatched(display.Self, l) {
			return l, true
		}
	}
	return language.English, false
}

func (p *MediaInfoParser) parseAudio(line string) {
	if sm := audioRegex.FindStringSubmatch(line); sm != nil {
		audioLine := sm[3]
		id := p.Info.CreateAudioTrack()
		if lang, ok := parseLanguage(audioLine); ok {
			p.Info.Audio[id].Language = lang.String()
		}
		p.Info.Audio[id].Voice = audioLine
	}
}

func parseSubtitlesTracks(s string) []string {
	var result []string
	t := ""
	inBraces := false
	for _, c := range s {
		if c == '(' || c == '[' {
			inBraces = true
		}
		if c == ')' || c == ']' {
			inBraces = false
		}
		if c == ',' && !inBraces {
			result = append(result, strings.TrimSpace(t))
			t = ""
		} else {
			t = t + string(c)
		}
	}

	if t != "" && strings.ToLower(t) != "отсутствуют" {
		result = append(result, strings.TrimSpace(t))
	}

	return result
}
func (p *MediaInfoParser) parseSubtitles(line string) {
	if sm := subtitlesRegex.FindStringSubmatch(line); sm != nil {
		subLine := sm[1]
		tracks := parseSubtitlesTracks(subLine)
		for _, s := range tracks {
			id := p.Info.CreateSubtitleTrack()
			if lang, ok := parseLanguage(s); ok {
				p.Info.Subtitle[id].Language = lang.String()
			} else {
				p.Info.Subtitle[id].Language = s
			}
		}
	}
}
