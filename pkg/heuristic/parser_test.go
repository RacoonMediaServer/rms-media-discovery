package heuristic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseName(t *testing.T) {
	type testCase struct {
		input  string
		output tokenList
	}

	testCases := []testCase{
		{
			input: "Стражи Галактики.2014.UHD.BDRip.2160p",
			output: tokenList{
				{Text: "стражи"},
				{Text: "галактики"},
				{Text: "2014"},
				{Text: "uhd"},
				{Text: "bdrip"},
				{Text: "2160p"},
			},
		},
		{
			input: "Люди Икс (BDRip)",
			output: tokenList{
				{Text: "люди"},
				{Text: "икс"},
				{Text: "bdrip", InBraces: true},
			},
		},
		{
			input: "Люди Икс [BDRip]",
			output: tokenList{
				{Text: "люди"},
				{Text: "икс"},
				{Text: "bdrip", InBraces: true},
			},
		},
		{
			input: "Люди Икс [BDRip 1080p",
			output: tokenList{
				{Text: "люди"},
				{Text: "икс"},
				{Text: "bdrip", InBraces: true},
				{Text: "1080p", InBraces: true},
			},
		},

		{
			input: "Люди Икс ((BDRip) 1080p) 2",
			output: tokenList{
				{Text: "люди"},
				{Text: "икс"},
				{Text: "bdrip", InBraces: true},
				{Text: "1080p", InBraces: true},
				{Text: "2"},
			},
		},
	}

	for i, tc := range testCases {
		actual := parse(tc.input)
		assert.Equal(t, tc.output, actual, "Test %d failed", i)
	}
}
