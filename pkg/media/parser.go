package media

import (
	"bufio"
	"strconv"
	"strings"
)

func ParseInfo(content string) *Info {
	type stream struct {
		params map[string]string
	}
	var streams []*stream
	var cur *stream

	mi := &Info{}

	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		text := scanner.Text()
		key, value, ok := strings.Cut(text, ": ")
		key = strings.TrimRight(key, " ")
		key = strings.TrimRight(key, string(rune(0xa0)))
		key = strings.TrimRight(key, "\t")
		if !ok {
			continue
		}
		if key == "ID" {
			cur = new(stream)
			cur.params = make(map[string]string)
			streams = append(streams, cur)
		}
		if cur != nil {
			cur.params[key] = value
		}
		if cur == nil && key == "Format" {
			mi.Format = value
		}
	}

	if len(streams) == 0 {
		return nil
	}

	for _, s := range streams {
		if _, ok := s.params["Width"]; ok {
			mi.Video = append(mi.Video, fillVideoTrack(s.params))
		} else if _, ok := s.params["Sampling rate"]; ok {
			mi.Audio = append(mi.Audio, fillAudioTrack(s.params))
		} else {
			mi.Subtitle = append(mi.Subtitle, fillSubtitleTrack(s.params))
		}
	}

	return mi
}

func fillVideoTrack(params map[string]string) VideoTrack {
	v := VideoTrack{}
	normalize := func(s string) string {
		s = strings.TrimSuffix(s, " pixels")
		return strings.ReplaceAll(s, " ", "")
	}
	if val, ok := params["Width"]; ok {
		conv, _ := strconv.ParseInt(normalize(val), 10, 32)
		v.Width = int(conv)
	}

	if val, ok := params["Height"]; ok {
		conv, _ := strconv.ParseInt(normalize(val), 10, 32)
		v.Height = int(conv)
	}

	if val, ok := params["Codec ID"]; ok {
		v.Codec = val
	}

	if val, ok := params["Display aspect ratio"]; ok {
		v.AspectRatio = val
	}

	return v
}

func fillAudioTrack(params map[string]string) AudioTrack {
	a := AudioTrack{}

	if val, ok := params["Codec ID"]; ok {
		a.Codec = val
	}

	if val, ok := params["Title"]; ok {
		a.Voice = val
	}

	if val, ok := params["Language"]; ok {
		a.Language = val
	}

	return a
}

func fillSubtitleTrack(params map[string]string) SubtitleTrack {
	s := SubtitleTrack{}

	if val, ok := params["Codec ID"]; ok {
		s.Codec = val
	}

	if val, ok := params["Language"]; ok {
		s.Language = val
	}

	return s
}
