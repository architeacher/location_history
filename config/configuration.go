package config

import "time"

type (
	// Configuration type wraps configuration data.
	Configuration struct {
		Port                           string
		TTL, WriteTimeout, ReadTimeout time.Duration
		BufferSize                     uint
	}
)
