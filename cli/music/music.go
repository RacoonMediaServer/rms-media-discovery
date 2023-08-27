package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/provider/deezer"
)

func main() {
	q := flag.String("query", "", "search query")
	flag.Parse()

	provider := deezer.NewProvider()
	results, err := provider.SearchMusic(context.Background(), *q, 10)
	if err != nil {
		panic(err)
	}
	for _, res := range results {
		if res.IsArtist() {
			fmt.Printf("Artist: %+v\n", res.AsArtist())
		} else if res.IsAlbum() {
			fmt.Printf("Album: %+v\n", res.AsAlbum())
		} else if res.IsTrack() {
			fmt.Printf("Track: %+v\n", res.AsTrack())
		}
	}
}
