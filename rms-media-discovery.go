package main

import (
	"flag"
	"github.com/RacoonMediaServer/rms-media-discovery/internal/db"
	"github.com/RacoonMediaServer/rms-media-discovery/internal/server"
	"github.com/RacoonMediaServer/rms-media-discovery/internal/service/accounts"
	"github.com/RacoonMediaServer/rms-media-discovery/internal/service/movies"
	"github.com/RacoonMediaServer/rms-media-discovery/internal/service/torrents"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/navigator"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/pipeline"
	"github.com/RacoonMediaServer/rms-packages/pkg/service/servicemgr"
	"github.com/apex/log"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go-micro.dev/v4"
	"net/http"
)

const version = "1.1.8"

func main() {
	log.Infof("rms-media-discovery v%s", version)
	host := flag.String("host", "127.0.0.1", "Server IP address")
	port := flag.Int("port", 8080, "Server port")
	dbString := flag.String("db", "mongodb://localhost:27017", "MongoDB connection string")
	verbose := flag.Bool("verbose", false, "Verbose mode")
	flag.Parse()

	log.Info("Headless browser engine initializing...")
	if err := navigator.Initialize(); err != nil {
		log.Fatalf("Failed: %s", err)
	}

	if *verbose {
		log.SetLevel(log.DebugLevel)
		navigator.SetSettings(navigator.Settings{StoreDumpOnError: true})
	}

	db, err := db.Connect(*dbString)
	if err != nil {
		log.Fatalf("Connect to database failed: %s", err)
	}
	log.Info("Connected to MongoDB")

	service := micro.NewService(micro.Name("rms-media-discovery"))

	srv := server.Server{}
	srv.Users = servicemgr.NewServiceFactory(service).NewUsers()
	srv.Accounts = accounts.New(db)
	srv.Movies = movies.New(srv.Accounts)
	srv.Torrents = torrents.New(srv.Accounts)

	if err := srv.Accounts.Initialize(); err != nil {
		log.Fatalf("Initialize accounts service failed: %+s", err)
	}

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		if err := http.ListenAndServe(":2112", nil); err != nil {
			log.Fatalf("Cannot bind monitoring endpoint: %s", err)
		}
	}()

	if err := srv.ListenAndServer(*host, *port); err != nil {
		log.Fatalf("Cannot start web server: %+s", err)
	}

	pipeline.Stop()
}
