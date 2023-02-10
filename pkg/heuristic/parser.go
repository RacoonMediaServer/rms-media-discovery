package heuristic

import (
	"unicode"
)

func parse(text string) tokenList {
	tokens := tokenList{}
	t := token{}
	braces := 0

	for _, ch := range text {
		if !unicode.IsLetter(ch) && !unicode.IsDigit(ch) {
			if !t.IsEmpty() {
				tokens.Push(t)
				t.Clear()
			}
			if ch == '(' || ch == '[' {
				braces++
				t.InBraces = true
			} else if (ch == ')' || ch == ']') && braces > 0 {
				braces--
				if braces == 0 {
					t.InBraces = false
				}
			} else if ch == '/' {
				t.SeqStart = true
			}

		} else {
			t.Push(ch)
		}
	}

	if !t.IsEmpty() {
		tokens.Push(t)
	}

	return tokens
}
