package unit

import (
	"github.com/w-h-a/pkg/telemetry/log"
	memorylog "github.com/w-h-a/pkg/telemetry/log/memory"
	"github.com/w-h-a/pkg/utils/memoryutils"
)

func SetLogger() {
	logBuffer := memoryutils.NewBuffer()

	logger := memorylog.NewLog(
		log.LogWithPrefix("export"),
		memorylog.LogWithBuffer(logBuffer),
	)

	log.SetLogger(logger)
}

func Bool(v bool) *bool {
	return &v
}

func Float64(v float64) *float64 {
	return &v
}

func Int(v int) *int {
	return &v
}

func String(v string) *string {
	return &v
}
