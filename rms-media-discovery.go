package main

import (
	"flag"

	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/server"
	"github.com/apex/log"
)

func main() {
	host := flag.String("host", "127.0.0.1", "Server IP address")
	port := flag.Int("port", 8080, "Server port")
	flag.Parse()

	srv := server.Server{}
	if err := srv.ListenAndServer(*host, *port); err != nil {
		log.Fatalf("Cannot start web server: %+s", err)
	}
}
