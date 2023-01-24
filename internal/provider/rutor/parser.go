package rutor

import (
	"regexp"
	"strconv"
)

var extractSizeExpr = regexp.MustCompile(`(\d+(.\d+)?).(MB|GB)`)

func parseTorrentSize(text string) uint64 {
	matches := extractSizeExpr.FindStringSubmatch(text)
	if matches != nil {
		result, err := strconv.ParseFloat(matches[1], 32)
		if err != nil {
			return 0
		}
		if matches[3] == "GB" {
			result *= 1024.
		}
		return uint64(result)
	}

	return 0
}
