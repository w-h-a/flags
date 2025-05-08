package config

import (
	"os"
	"strings"
	"sync"
)

var (
	instance *config
	once     sync.Once
)

type config struct {
	env             string
	name            string
	version         string
	httpAddress     string
	tracesAddress   string
	metricsAddress  string
	flagFormat      string
	fileClient      string
	fileClientDir   string
	fileClientFiles []string
	fileClientToken string
	reportClient    string
	reportClientDir string
	messageClient   string
}

func New() {
	once.Do(func() {
		instance = &config{
			env:             "dev",
			name:            "flags",
			version:         "0.1.0-alpha.0",
			httpAddress:     ":0",
			tracesAddress:   "localhost:4318",
			metricsAddress:  "localhost:4318",
			flagFormat:      "yaml",
			fileClient:      "local",
			fileClientDir:   ".",
			fileClientFiles: []string{"/flags.yaml"},
			fileClientToken: "",
			reportClient:    "local",
			reportClientDir: "/tmp",
			messageClient:   "local",
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

		metricsAddress := os.Getenv("METRICS_ADDRESS")
		if len(metricsAddress) > 0 {
			instance.metricsAddress = metricsAddress
		}

		flagFormat := os.Getenv("FLAG_FORMAT")
		if len(flagFormat) > 0 {
			instance.flagFormat = flagFormat
		}

		fileClient := os.Getenv("FILE_CLIENT")
		if len(fileClient) > 0 {
			instance.fileClient = fileClient
		}

		fileClientDir := os.Getenv("FILE_CLIENT_DIR")
		if len(fileClientDir) > 0 {
			instance.fileClientDir = fileClientDir
		}

		fileClientFiles := os.Getenv("FILE_CLIENT_FILES")
		if len(fileClientFiles) > 0 {
			instance.fileClientFiles = append(instance.fileClientFiles, strings.Split(fileClientFiles, ",")...)
		}

		fileClientToken := os.Getenv("FILE_CLIENT_TOKEN")
		if len(fileClientToken) > 0 {
			instance.fileClientToken = fileClientToken
		}

		reportClient := os.Getenv("REPORT_CLIENT")
		if len(reportClient) > 0 {
			instance.reportClient = reportClient
		}

		reportClientDir := os.Getenv("REPORT_CLIENT_DIR")
		if len(reportClientDir) > 0 {
			instance.reportClientDir = reportClientDir
		}

		messageClient := os.Getenv("MESSAGE_CLIENT")
		if len(messageClient) > 0 {
			instance.messageClient = messageClient
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

func MetricsAddress() string {
	if instance == nil {
		return ""
	}

	return instance.metricsAddress
}

func FlagFormat() string {
	if instance == nil {
		return ""
	}

	return instance.flagFormat
}

func FileClient() string {
	if instance == nil {
		return ""
	}

	return instance.fileClient
}

func FileClientDir() string {
	if instance == nil {
		return ""
	}

	return instance.fileClientDir
}

func FileClientFiles() []string {
	if instance == nil {
		return []string{}
	}

	return instance.fileClientFiles
}

func FileClientToken() string {
	if instance == nil {
		return ""
	}

	return instance.fileClientToken
}

func ReportClient() string {
	if instance == nil {
		return ""
	}

	return instance.reportClient
}

func ReportClientDir() string {
	if instance == nil {
		return ""
	}

	return instance.reportClientDir
}

func MessageClient() string {
	if instance == nil {
		return ""
	}

	return instance.messageClient
}

// used for test purposes only
func Reset() {
	instance = &config{
		env:             "dev",
		name:            "flags",
		version:         "0.1.0-alpha.0",
		httpAddress:     ":0",
		tracesAddress:   "localhost:4318",
		metricsAddress:  "localhost:4318",
		flagFormat:      "yaml",
		fileClient:      "local",
		fileClientDir:   ".",
		fileClientFiles: []string{"/flags.yaml"},
		fileClientToken: "",
		reportClient:    "local",
		reportClientDir: "/tmp",
		messageClient:   "local",
	}

	once = sync.Once{}
}
