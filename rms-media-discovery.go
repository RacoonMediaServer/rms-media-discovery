package main

import (
	"fmt"
	"github.com/RacoonMediaServer/rms-media-discovery/internal/mocks"
	"github.com/RacoonMediaServer/rms-media-discovery/internal/service/music"
	"net/http"

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
)

var Version = "0.0.0"

func main() {
	log.Infof("rms-media-discovery %s", Version)
	defer log.Info("DONE.")

	log.Info("Headless browser engine initializing...")
	if err := navigator.Initialize(); err != nil {
		log.Fatalf("Failed: %s", err)
	}

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

	cfg := config.Config()

	if useDebug || cfg.Debug.Verbose {
		log.SetLevel(log.DebugLevel)
		navigator.SetSettings(navigator.Settings{StoreDumpOnError: true})
	}

	database, err := db.Connect(cfg.Database)
	if err != nil {
		log.Fatalf("Connect to database failed: %s", err)
	}
	log.Info("Connected to database")

	accountsService := accounts.New(database)
	if err = accountsService.Initialize(); err != nil {
		log.Fatalf("Initialize accounts service failed: %s", err)
	}
	usersService := servicemgr.NewServiceFactory(service).NewUsers()
	if cfg.DisableAccessControl {
		usersService = mocks.NewMockUsersAllAllowed()
	}
	torrentService := torrents.New(accountsService)

	srv := server.Server{
		Movies:   movies.New(accountsService),
		Music:    music.New(),
		Torrents: torrentService,
		Users:    usersService,
		Accounts: accountsService,
	}

	if cfg.Debug.Monitor.Enabled {
		go func() {
			http.Handle("/metrics", promhttp.Handler())
			if err := http.ListenAndServe(fmt.Sprintf("%s:%d", cfg.Debug.Monitor.Host, cfg.Debug.Monitor.Port), nil); err != nil {
				log.Fatalf("Cannot bind monitoring endpoint: %s", err)
			}
		}()
	}

	if err := srv.ListenAndServer(cfg.Http.Host, cfg.Http.Port); err != nil {
		log.Fatalf("Cannot start web server: %+s", err)
	}

	pipeline.Stop()
	torrentService.Stop()
}
