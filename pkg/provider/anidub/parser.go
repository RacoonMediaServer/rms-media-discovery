package anidub

import (
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/scraper"
	"github.com/gocolly/colly/v2"
)

func loginChecker(isLogged *bool) scraper.HTMLCallback {
	return func(e *colly.HTMLElement, userData interface{}) {

	}
}
