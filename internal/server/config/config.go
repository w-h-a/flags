package config

import (
	"os"
	"strconv"
	"strings"
	"sync"
)

var (
	instance *config
	once     sync.Once
)

type config struct {
	env                 string
	region              string
	name                string
	version             string
	httpAddress         string
	apiKeys             map[string]bool
	logsExporter        string
	logsAddress         string
	logsUrlPath         string
	logsAPIToken        string
	tracesExporter      string
	tracesAddress       string
	metricsExporter     string
	metricsAddress      string
	flagFormat          string
	writeClient         string
	writeClientLocation string
	writeClientToken    string
	readClient          string
	readClientLocation  string
	readClientToken     string
	readInterval        int
	exportReports       bool
	exportClient        string
	exportClientDir     string
	exportInterval      int
	notifyClient        string
	notifyURL           string
}

func New() {
	once.Do(func() {
		instance = &config{
			env:                 "dev",
			region:              "local",
			name:                "flags",
			version:             "0.1.0-alpha.0",
			httpAddress:         ":0",
			apiKeys:             map[string]bool{},
			logsExporter:        "stdout",
			logsAddress:         "",
			logsUrlPath:         "",
			logsAPIToken:        "",
			tracesExporter:      "otlp",
			tracesAddress:       "localhost:4318",
			metricsExporter:     "otlp",
			metricsAddress:      "localhost:4318",
			flagFormat:          "yaml",
			writeClient:         "noop",
			writeClientLocation: "noop",
			writeClientToken:    "",
			readClient:          "local",
			readClientLocation:  "./flags.yaml",
			readClientToken:     "",
			readInterval:        60,
			exportReports:       false,
			exportClient:        "local",
			exportClientDir:     "/tmp",
			exportInterval:      120,
			notifyClient:        "local",
			notifyURL:           "",
		}

		env := os.Getenv("ENV")
		if len(env) > 0 {
			instance.env = env
		}

		region := os.Getenv("AWS_REGION")
		if len(region) > 0 {
			instance.region = region
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

		apiKeys := os.Getenv("API_KEYS")
		if len(apiKeys) > 0 {
			keys := strings.Split(apiKeys, ",")
			for _, k := range keys {
				instance.apiKeys[k] = true
			}
		}

		logsExporter := os.Getenv("LOGS_EXPORTER")
		if len(logsExporter) > 0 {
			instance.logsExporter = logsExporter
		}

		logsAddress := os.Getenv("LOGS_ADDRESS")
		if len(logsAddress) > 0 {
			instance.logsAddress = logsAddress
		}

		logsUrlPath := os.Getenv("LOGS_URL_PATH")
		if len(logsUrlPath) > 0 {
			instance.logsUrlPath = logsUrlPath
		}

		logsAPIToken := os.Getenv("LOGS_API_TOKEN")
		if len(logsAPIToken) > 0 {
			instance.logsAPIToken = logsAPIToken
		}

		tracesExporter := os.Getenv("TRACES_EXPORTER")
		if len(tracesExporter) > 0 {
			instance.tracesExporter = tracesExporter
		}

		tracesAddress := os.Getenv("TRACES_ADDRESS")
		if len(tracesAddress) > 0 {
			instance.tracesAddress = tracesAddress
		}

		metricsExporter := os.Getenv("METRICS_EXPORTER")
		if len(metricsExporter) > 0 {
			instance.metricsExporter = metricsExporter
		}

		metricsAddress := os.Getenv("METRICS_ADDRESS")
		if len(metricsAddress) > 0 {
			instance.metricsAddress = metricsAddress
		}

		flagFormat := os.Getenv("FLAG_FORMAT")
		if len(flagFormat) > 0 {
			instance.flagFormat = flagFormat
		}

		writeClient := os.Getenv("WRITE_CLIENT")
		if len(writeClient) > 0 {
			instance.writeClient = writeClient
		}

		writeClientLocation := os.Getenv("WRITE_CLIENT_LOCATION")
		if len(writeClientLocation) > 0 {
			instance.writeClientLocation = writeClientLocation
		}

		writeClientToken := os.Getenv("WRITE_CLIENT_TOKEN")
		if len(writeClientToken) > 0 {
			instance.writeClientToken = writeClientToken
		}

		readClient := os.Getenv("READ_CLIENT")
		if len(readClient) > 0 {
			instance.readClient = readClient
		}

		readClientLocation := os.Getenv("READ_CLIENT_LOCATION")
		if len(readClientLocation) > 0 {
			instance.readClientLocation = readClientLocation
		}

		readClientToken := os.Getenv("READ_CLIENT_TOKEN")
		if len(readClientToken) > 0 {
			instance.readClientToken = readClientToken
		}

		readInterval := os.Getenv("READ_INTERVAL")
		if len(readInterval) > 0 {
			if interval, err := strconv.Atoi(readInterval); err == nil && interval >= 1 {
				instance.readInterval = interval
			}
		}

		exportReports := os.Getenv("EXPORT_REPORTS")
		if len(exportReports) > 0 {
			if exportReports == "true" {
				instance.exportReports = true

				exportClient := os.Getenv("EXPORT_CLIENT")
				if len(exportClient) > 0 {
					instance.exportClient = exportClient
				}

				exportClientDir := os.Getenv("EXPORT_CLIENT_DIR")
				if len(exportClientDir) > 0 {
					instance.exportClientDir = exportClientDir
				}

				exportInterval := os.Getenv("EXPORT_INTERVAL")
				if len(exportInterval) > 0 {
					if interval, err := strconv.Atoi(exportInterval); err == nil && interval >= 1 {
						instance.exportInterval = interval
					}
				}
			}
		}

		notifyClient := os.Getenv("NOTIFY_CLIENT")
		if len(notifyClient) > 0 {
			instance.notifyClient = notifyClient
		}

		notifyURL := os.Getenv("NOTIFY_URL")
		if len(notifyURL) > 0 {
			instance.notifyURL = notifyURL
		}
	})
}

func Env() string {
	if instance == nil {
		return ""
	}

	return instance.env
}

func Region() string {
	if instance == nil {
		return ""
	}

	return instance.region
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

func CheckAPIKey(key string) bool {
	if instance == nil {
		return false
	}

	_, ok := instance.apiKeys[key]
	return ok
}

func LogsExporter() string {
	if instance == nil {
		return ""
	}

	return instance.logsExporter
}

func LogsAddress() string {
	if instance == nil {
		return ""
	}

	return instance.logsAddress
}

func LogsUrlPath() string {
	if instance == nil {
		return ""
	}

	return instance.logsUrlPath
}

func LogsAPIToken() string {
	if instance == nil {
		return ""
	}

	return instance.logsAPIToken
}

func TracesExporter() string {
	if instance == nil {
		return ""
	}

	return instance.tracesExporter
}

func TracesAddress() string {
	if instance == nil {
		return ""
	}

	return instance.tracesAddress
}

func MetricsExporter() string {
	if instance == nil {
		return ""
	}

	return instance.metricsExporter
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

func WriteClient() string {
	if instance == nil {
		return ""
	}

	return instance.writeClient
}

func WriteClientLocation() string {
	if instance == nil {
		return ""
	}

	return instance.writeClientLocation
}

func WriteClientToken() string {
	if instance == nil {
		return ""
	}

	return instance.writeClientToken
}

func ReadClient() string {
	if instance == nil {
		return ""
	}

	return instance.readClient
}

func ReadClientLocation() string {
	if instance == nil {
		return ""
	}

	return instance.readClientLocation
}

func ReadClientToken() string {
	if instance == nil {
		return ""
	}

	return instance.readClientToken
}

func ReadInterval() int {
	if instance == nil {
		return 0
	}

	return instance.readInterval
}

func ExportReports() bool {
	if instance == nil {
		return false
	}

	return instance.exportReports
}

func ExportClient() string {
	if instance == nil {
		return ""
	}

	return instance.exportClient
}

func ExportClientDir() string {
	if instance == nil {
		return ""
	}

	return instance.exportClientDir
}

func ExportInterval() int {
	if instance == nil {
		return 0
	}

	return instance.exportInterval
}

func NotifyClient() string {
	if instance == nil {
		return ""
	}

	return instance.notifyClient
}

func NotifyURL() string {
	if instance == nil {
		return ""
	}

	return instance.notifyURL
}

// used for test purposes only
func Reset() {
	instance = &config{
		env:                 "dev",
		region:              "local",
		name:                "flags",
		version:             "0.1.0-alpha.0",
		httpAddress:         ":0",
		apiKeys:             map[string]bool{},
		tracesAddress:       "localhost:4318",
		metricsAddress:      "localhost:4318",
		flagFormat:          "yaml",
		writeClient:         "noop",
		writeClientLocation: "noop",
		writeClientToken:    "",
		readClient:          "local",
		readClientLocation:  "./flags.yaml",
		readClientToken:     "",
		readInterval:        60,
		exportReports:       false,
		exportClient:        "local",
		exportClientDir:     "/tmp",
		exportInterval:      120,
		notifyClient:        "local",
		notifyURL:           "",
	}

	once = sync.Once{}
}
