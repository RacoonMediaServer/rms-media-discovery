package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/model"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/provider"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/provider/imdb"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/provider/kinopoisk"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/provider/tmdb"
)

type accessor struct {
	key string
}

func (a accessor) GetCredentials(serviceId string) (model.Credentials, error) {
	//TODO implement me
	panic("implement me")
}

func (a accessor) GetApiKey(serviceId string) (model.ApiKey, error) {
	return model.ApiKey{AccountId: serviceId, Key: a.key}, nil
}

func (a accessor) MarkUnaccesible(accountId string) {
	//TODO implement me
	panic("implement me")
}

func main() {
	var a accessor

	query := flag.String("query", "", "query for search")
	id := flag.String("id", "", "movie id")
	key := flag.String("key", "", "API key")
	prov := flag.String("provider", "tmdb", "provider")
	flag.Parse()

	if *query == "" && *id == "" {
		panic("query or id must be set")
	}
	if *prov == "" {
		panic("provider muse be set")
	}

	a.key = *key

	var source provider.MovieInfoProvider
	switch *prov {
	case "kinopoisk":
		source = kinopoisk.NewProvider(a)
	case "imdb":
		source = imdb.NewProvider(a)
	case "tmdb":
		source = tmdb.NewProvider(a)
	default:
		panic(fmt.Sprintf("%s provider not implemented", *prov))
	}

	if *query != "" {
		results, err := source.SearchMovies(context.Background(), *query, 5)
		if err != nil {
			panic(err)
		}
		for _, r := range results {
			fmt.Println(r)
		}
	}

	if *id != "" {
		result, err := source.GetMovieInfo(context.Background(), *id)
		if err != nil {
			panic(err)
		}
		fmt.Println(result)
	}
}
