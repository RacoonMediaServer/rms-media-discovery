package main

import (
	"flag"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/pipeline"

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
	verbose := flag.Bool("verbose", false, "Verbose mode")
	flag.Parse()

	if *verbose {
		log.SetLevel(log.DebugLevel)
	}

	db, err := db.Connect(*dbString)
	if err != nil {
		log.Fatalf("Connect to database failed: %+w", err)
	}
	log.Info("Connected to MongoDB")

	srv := server.Server{}
	srv.Users = users.New(db)
	srv.Accounts = accounts.New(db)
	srv.Movies = movies.New(srv.Accounts)
	srv.Torrents = torrents.New()

	if err := srv.Users.Initialize(); err != nil {
		log.Fatalf("Initialize users service failed: %+s", err)
	}

	if err := srv.Accounts.Initialize(); err != nil {
		log.Fatalf("Initialize accounts service failed: %+s", err)
	}

	if err := srv.ListenAndServer(*host, *port); err != nil {
		log.Fatalf("Cannot start web server: %+s", err)
	}

	pipeline.Stop()
}
