package heuristic

import (
	"bufio"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/pkg/media"
	"strconv"
	"strings"
)

type paramMap map[string]string

func (m paramMap) Get(keys ...string) (string, bool) {
	for i := range keys {
		if v, ok := m[keys[i]]; ok {
			return v, true
		}
	}

	return "", false
}

func parseMediaInfo(content string) *media.Info {
	type stream struct {
		params paramMap
	}
	var streams []*stream
	var cur *stream

	mi := &media.Info{}

	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		text := scanner.Text()
		key, value, ok := strings.Cut(text, ": ")
		key = strings.TrimRight(key, " ")
		key = strings.TrimRight(key, "\t")
		if !ok {
			continue
		}
		if key == "ID" || key == "Идентификатор" {
			cur = new(stream)
			cur.params = make(paramMap)
			streams = append(streams, cur)
		}
		if cur != nil {
			cur.params[key] = value
		}
		if cur == nil && (key == "Format" || key == "Формат") {
			mi.Format = value
		}
	}

	if len(streams) == 0 {
		return mi
	}

	for _, s := range streams {
		if _, ok := s.params.Get("Width", "Ширина"); ok {
			mi.Video = append(mi.Video, fillVideoTrack(s.params))
		} else if _, ok := s.params.Get("Sampling rate", "Частота"); ok {
			mi.Audio = append(mi.Audio, fillAudioTrack(s.params))
		} else {
			mi.Subtitle = append(mi.Subtitle, fillSubtitleTrack(s.params))
		}
	}

	return mi
}

func fillVideoTrack(params paramMap) media.VideoTrack {
	v := media.VideoTrack{}
	normalize := func(s string) string {
		s = strings.TrimSuffix(s, " pixels")
		s = strings.TrimSuffix(s, " пикс.")
		return strings.ReplaceAll(s, " ", "")
	}
	if val, ok := params.Get("Width", "Ширина"); ok {
		conv, _ := strconv.ParseInt(normalize(val), 10, 32)
		v.Width = int(conv)
	}

	if val, ok := params.Get("Height", "Высота"); ok {
		conv, _ := strconv.ParseInt(normalize(val), 10, 32)
		v.Height = int(conv)
	}

	if val, ok := params.Get("Codec ID", "Идентификатор кодека"); ok {
		v.Codec = val
	}

	if val, ok := params.Get("Display aspect ratio", "Соотношение кадра"); ok {
		v.AspectRatio = val
	}

	return v
}

func fillAudioTrack(params paramMap) media.AudioTrack {
	a := media.AudioTrack{}

	if val, ok := params.Get("Codec ID", "Идентификатор кодека"); ok {
		a.Codec = val
	}

	if val, ok := params.Get("Title", "Заголовок"); ok {
		a.Voice = val
	}

	if val, ok := params.Get("Language", "Язык"); ok {
		if l, ok := parseLanguage(val); ok {
			a.Language = l.String()
		} else {
			a.Language = val
		}
	}

	return a
}

func fillSubtitleTrack(params paramMap) media.SubtitleTrack {
	s := media.SubtitleTrack{}

	if val, ok := params.Get("Codec ID", "Идентификатор кодека"); ok {
		s.Codec = val
	}

	if val, ok := params.Get("Language", "Язык"); ok {
		if l, ok := parseLanguage(val); ok {
			s.Language = l.String()
		} else {
			s.Language = val
		}
	}

	return s
}
