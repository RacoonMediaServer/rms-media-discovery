package heuristic

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNormalize(t *testing.T) {
	text := "AC/DC - Дискография 1974-2014 (131 релиз, включая бокс сеты), MP3, 320 kbps"
	expected := "ac dc дискография 1974 2014 131 релиз включая бокс сеты mp3 320 kbps"
	assert.Equal(t, expected, Normalize(text))
}

func TestNormalizeWithoutBraces(t *testing.T) {
	text := "AC/DC - Дискография 1974-2014 (131 релиз, включая бокс сеты), MP3, 320 kbps"
	expected := "ac dc дискография 1974 2014 mp3 320 kbps"
	assert.Equal(t, expected, NormalizeWithoutBraces(text))
}
