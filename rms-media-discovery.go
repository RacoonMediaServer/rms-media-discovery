package main

import (
	"flag"

	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/server"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/service/admin"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/service/movies"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/service/torrents"
	"github.com/apex/log"
)

func main() {
	host := flag.String("host", "127.0.0.1", "Server IP address")
	port := flag.Int("port", 8080, "Server port")
	flag.Parse()

	srv := server.Server{}
	srv.Admin = admin.New()
	srv.Movies = movies.New()
	srv.Torrents = torrents.New()

	if err := srv.ListenAndServer(*host, *port); err != nil {
		log.Fatalf("Cannot start web server: %+s", err)
	}
}
