package config

import (
	"os"
	"sync"
)

var (
	instance *config
	once     sync.Once
)

type config struct {
	env           string
	name          string
	version       string
	httpAddress   string
	tracesAddress string
	// metricsAddress string
}

func New() {
	once.Do(func() {
		instance = &config{
			env:           "dev",
			name:          "flags",
			version:       "0.1.0-alpha.0",
			httpAddress:   ":0",
			tracesAddress: "localhost:4318",
		}

		env := os.Getenv("ENV")
		if len(env) > 0 {
			instance.env = env
		}

		name := os.Getenv("NAME")
		if len(name) > 0 {
			instance.name = name
		}

		version := os.Getenv("VERSION")
		if len(version) > 0 {
			instance.version = version
		}

		httpAddress := os.Getenv("HTTP_ADDRESS")
		if len(httpAddress) > 0 {
			instance.httpAddress = httpAddress
		}

		tracesAddress := os.Getenv("TRACES_ADDRESS")
		if len(tracesAddress) > 0 {
			instance.tracesAddress = tracesAddress
		}
	})
}

func Env() string {
	if instance == nil {
		return ""
	}

	return instance.env
}

func Name() string {
	if instance == nil {
		return ""
	}

	return instance.name
}

func Version() string {
	if instance == nil {
		return ""
	}

	return instance.version
}

func HttpAddress() string {
	if instance == nil {
		return ""
	}

	return instance.httpAddress
}

func TracesAddress() string {
	if instance == nil {
		return ""
	}

	return instance.tracesAddress
}
