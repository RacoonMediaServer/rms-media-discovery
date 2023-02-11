package heuristic

import (
	"unicode"
	"unicode/utf8"
)

type token struct {
	Text     string
	InBraces bool
	SeqStart bool
}

type tokenList []token

func (t *token) Push(r rune) {
	t.Text += string(unicode.ToLower(r))
}

func (t token) IsEmpty() bool {
	return t.Text == ""
}

func (t token) IsDigital() bool {
	for _, ch := range t.Text {
		if !unicode.IsDigit(ch) {
			return false
		}
	}

	return true
}

func (t token) String() string {
	if utf8.RuneCountInString(t.Text) <= 2 {
		return t.Text
	}
	result := ""
	for i, r := range t.Text {
		if i == 0 {
			result += string(unicode.ToUpper(r))
			continue
		}
		result += string(r)
	}
	return result
}

func (t *tokenList) Push(tok token) {
	*t = append(*t, tok)
}

func (t tokenList) Find(m match) int {
	for i := range t {
		if m.Match(t[i]) {
			return i
		}
	}

	return -1
}

func (t tokenList) Remove(i int) tokenList {
	slice := make([]token, 0, len(t))
	slice = append(slice, t[:i]...)
	slice = append(slice, t[i+1:]...)
	return slice
}

func (t *token) Clear() {
	t.Text = ""
	t.SeqStart = false
}

func (t tokenList) RemoveIf(m match) tokenList {
	slice := tokenList{}
	for i := range t {
		if !m.Match(t[i]) {
			slice = append(slice, t[i])
		}
	}
	return slice
}

func (t tokenList) FindAll(m match) []int {
	result := make([]int, 0, len(t))
	for i := range t {
		if m.Match(t[i]) {
			result = append(result, i)
		}
	}
	return result
}

func (t tokenList) String() string {
	s := ""
	for i := range t {
		s += t[i].String() + " "
	}
	if len(s) != 0 {
		s = s[:len(s)-1]
	}
	return s
}
