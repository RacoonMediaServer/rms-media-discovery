package config

import "github.com/RacoonMediaServer/rms-packages/pkg/configuration"

// Configuration represents entire service configuration
type Configuration struct {
	Database             configuration.Database
	Http                 configuration.Http
	Debug                configuration.Debug
	DisableAccessControl bool
	MirrorBaseURL        string

	// PreferIPv6 - исходящие HTTP-запросы сначала пробуют IPv6 (RMS_PREFERIPV6=true).
	// Нужно, когда cookies (cf_clearance) получены с IPv6-адреса, а контейнер Docker из-за ULA-подсети выбирает IPv4
	PreferIPv6 bool
}

var config Configuration

// Load open and parses configuration file
func Load(configFilePath string) error {
	return configuration.Load(configFilePath, &config)
}

// Config returns loaded configuration
func Config() Configuration {
	return config
}
