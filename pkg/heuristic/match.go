package heuristic

import "regexp"

type match interface {
	Match(t token) bool
}

type wordMatch struct {
	Word string
}

type regexMatch struct {
	Exp *regexp.Regexp
}

type orMatch struct {
	Matches []match
}

type bracesMatch struct {
}

func (m wordMatch) Match(t token) bool {
	return t.Text == m.Word
}

func (m regexMatch) Match(t token) bool {
	return m.Exp.MatchString(t.Text)
}

func (m orMatch) Match(t token) bool {
	for _, ma := range m.Matches {
		if ma.Match(t) {
			return true
		}
	}

	return false
}

func (m bracesMatch) Match(t token) bool {
	return t.InBraces
}
