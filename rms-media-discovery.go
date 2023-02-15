package main

import (
	"github.com/RacoonMediaServer/rms-media-discovery/internal/config"
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
	"github.com/urfave/cli/v2"
	"go-micro.dev/v4"
	"net/http"
)

var Version = "0.0.0"

func main() {
	log.Infof("rms-media-discovery v%s", Version)
	defer log.Info("DONE.")

	useDebug := false

	service := micro.NewService(
		micro.Name("rms-media-discovery"),
		micro.Version(Version),
		micro.Flags(
			&cli.BoolFlag{
				Name:        "verbose",
				Aliases:     []string{"debug"},
				Usage:       "debug log level",
				Value:       false,
				Destination: &useDebug,
			},
		),
	)

	service.Init(
		micro.Action(func(context *cli.Context) error {
			configFile := "/etc/rms/rms-media-discovery.json"
			if context.IsSet("config") {
				configFile = context.String("config")
			}
			return config.Load(configFile)
		}),
	)

	if useDebug {
		log.SetLevel(log.DebugLevel)
		navigator.SetSettings(navigator.Settings{StoreDumpOnError: true})
	}

	log.Info("Headless browser engine initializing...")
	if err := navigator.Initialize(); err != nil {
		log.Fatalf("Failed: %s", err)
	}

	cfg := config.Config()

	database, err := db.Connect(cfg.Database)
	if err != nil {
		log.Fatalf("Connect to database failed: %s", err)
	}
	log.Info("Connected to MongoDB")

	srv := server.Server{}
	srv.Users = servicemgr.NewServiceFactory(service).NewUsers()
	srv.Accounts = accounts.New(database)
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

	if err := srv.ListenAndServer(cfg.Http.Host, cfg.Http.Port); err != nil {
		log.Fatalf("Cannot start web server: %+s", err)
	}

	pipeline.Stop()
}
