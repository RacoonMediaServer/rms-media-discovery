package main

import (
	"flag"

	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/db"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/server"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/service/accounts"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/service/movies"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/service/torrents"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/service/users"
	"github.com/apex/log"
)

const version = "0.0.1"

func main() {
	log.Infof("rms-media-discovery v%s", version)
	host := flag.String("host", "127.0.0.1", "Server IP address")
	port := flag.Int("port", 8080, "Server port")
	dbString := flag.String("db", "mongodb://localhost:27017", "MongoDB connection string")
	flag.Parse()

	_, err := db.Connect(*dbString)
	if err != nil {
		log.Fatalf("Connect to database failed: %+w", err)
	}

	srv := server.Server{}
	srv.Users = users.New()
	srv.Accounts = accounts.New()
	srv.Movies = movies.New()
	srv.Torrents = torrents.New()

	if err := srv.ListenAndServer(*host, *port); err != nil {
		log.Fatalf("Cannot start web server: %+s", err)
	}
}
