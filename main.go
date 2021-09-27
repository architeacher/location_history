package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/ahmedkamals/location_history/app"
	"github.com/ahmedkamals/location_history/config"
	"os"
	"strconv"
	"time"
)

const (
	defaultPort      string = ":8080"
	defaultTTL              = 60 * time.Second
	defaultQueueSize        = 100
)

func main() {
	config := parseFlags()
	logo()

	app.NewDispatcher(config).Dispatch(context.Background())
}

func logo() {
	fmt.Println(`
██       ██████   ██████  █████  ████████ ██  ██████  ███    ██     ██   ██ ██ ███████ ████████  ██████  ██████  ██    ██
██      ██    ██ ██      ██   ██    ██    ██ ██    ██ ████   ██     ██   ██ ██ ██         ██    ██    ██ ██   ██  ██  ██
██      ██    ██ ██      ███████    ██    ██ ██    ██ ██ ██  ██     ███████ ██ ███████    ██    ██    ██ ██████    ████
██      ██    ██ ██      ██   ██    ██    ██ ██    ██ ██  ██ ██     ██   ██ ██      ██    ██    ██    ██ ██   ██    ██
███████  ██████   ██████ ██   ██    ██    ██  ██████  ██   ████     ██   ██ ██ ███████    ██     ██████  ██   ██    ██
`)
}

func parseFlags() config.Configuration {
	bufferSize := flag.Uint("buffer_size", defaultQueueSize, "Max buffer size for a queue.")
	flag.Parse()

	port := os.Getenv("HISTORY_SERVER_LISTEN_ADDR")
	ttlInSeconds := os.Getenv("LOCATION_HISTORY_TTL_SECONDS")

	if port == "" {
		port = defaultPort
	}

	ttl := defaultTTL
	if ttlInSeconds != "" {
		ttlValue, _ := strconv.Atoi(ttlInSeconds)
		ttl = time.Second * time.Duration(ttlValue)
	}

	return config.Configuration{
		Port:       port,
		TTL:        ttl,
		BufferSize: *bufferSize,
	}
}
