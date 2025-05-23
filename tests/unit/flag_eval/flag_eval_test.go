package flageval

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/w-h-a/flags/internal/server"
	"github.com/w-h-a/flags/internal/server/clients/exporter"
	localexporter "github.com/w-h-a/flags/internal/server/clients/exporter/local"
	localnotifier "github.com/w-h-a/flags/internal/server/clients/notifier/local"
	"github.com/w-h-a/flags/internal/server/clients/reader"
	localreader "github.com/w-h-a/flags/internal/server/clients/reader/local"
	"github.com/w-h-a/flags/internal/server/clients/writer"
	"github.com/w-h-a/flags/internal/server/clients/writer/noop"
	"github.com/w-h-a/flags/internal/server/config"
	"github.com/w-h-a/pkg/telemetry/log"
	memorylog "github.com/w-h-a/pkg/telemetry/log/memory"
	"github.com/w-h-a/pkg/utils/memoryutils"
)

const (
	tok = "mytoken"
	dir = "../testdata"
)

func TestFlagEval_YAML(t *testing.T) {
	if len(os.Getenv("INTEGRATION")) > 0 {
		t.Log("SKIPPING UNIT TEST")
		return
	}

	// TODO: add tests when we read the req body
	type args struct {
		flagKey string
	}

	type want struct {
		httpCode int
		bodyFile string
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "bare-minimum-flag",
			args: args{
				flagKey: "bare-minimum-flag",
			},
			want: want{
				httpCode: http.StatusOK,
				bodyFile: "../testdata/flag_eval/bare_minimum_response.json",
			},
		},
		{
			name: "bare-minimum-flag-2",
			args: args{
				flagKey: "bare-minimum-flag-2",
			},
			want: want{
				httpCode: http.StatusOK,
				bodyFile: "../testdata/flag_eval/bare_minimum_response_2.json",
			},
		},
		{
			name: "get default value if flag is disabled",
			args: args{
				flagKey: "disabled-flag",
			},
			want: want{
				httpCode: http.StatusOK,
				bodyFile: "../testdata/flag_eval/disabled_flag_response.json",
			},
		},
		{
			name: "valid flag",
			args: args{
				flagKey: "allow-access",
			},
			want: want{
				httpCode: http.StatusOK,
				bodyFile: "../testdata/flag_eval/valid_response.json",
			},
		},
		{
			name: "request non-existant flag",
			args: args{
				flagKey: "not-found",
			},
			want: want{
				httpCode: http.StatusNotFound,
				bodyFile: "../testdata/flag_eval/flag_not_found_response.json",
			},
		},
	}

	for _, test := range tests {
		// env vars
		os.Setenv("API_KEYS", tok)
		os.Setenv("READ_CLIENT_LOCATION", dir+"/flags.yaml")

		// config
		config.New()

		// resource
		name := config.Name()

		// log
		logBuffer := memoryutils.NewBuffer()

		logger := memorylog.NewLog(
			log.LogWithPrefix(name),
			memorylog.LogWithBuffer(logBuffer),
		)

		log.SetLogger(logger)

		// traces

		// metrics

		// clients
		writeClient := noop.NewWriter(
			writer.WithLocation(config.WriteClientLocation()),
		)

		readClient := localreader.NewReader(
			reader.WithLocation(config.ReadClientLocation()),
		)

		exportClient := localexporter.NewExporter(
			exporter.WithDir(config.ExportClientDir()),
		)

		notifyClient := localnotifier.NewNotifier()

		// servers
		httpServer, _, exportService, notifyService, err := server.Factory(
			writeClient,
			readClient,
			exportClient,
			notifyClient,
		)
		require.NoError(t, err)

		t.Run(test.name, func(t *testing.T) {
			err = httpServer.Run()
			require.NoError(t, err)

			req, err := http.NewRequest(
				http.MethodPost,
				fmt.Sprintf("http://%s%s%s", httpServer.Options().Address, "/ofrep/v1/evaluate/flags/", test.args.flagKey),
				strings.NewReader(""),
			)
			require.NoError(t, err)

			req.Header.Set("content-type", "application/json")
			req.Header.Set("authorization", fmt.Sprintf("Bearer %s", tok))

			client := &http.Client{}

			rsp, err := client.Do(req)
			require.NoError(t, err)

			want, err := os.ReadFile(test.want.bodyFile)
			require.NoError(t, err)

			got, err := io.ReadAll(rsp.Body)
			require.NoError(t, err)

			require.Equal(t, string(want), string(got))

			require.Equal(t, test.want.httpCode, rsp.StatusCode)

			t.Cleanup(func() {
				rsp.Body.Close()
				exportService.Close()
				notifyService.Close()
				config.Reset()
			})
		})
	}
}

func TestFlagEval_JSON(t *testing.T) {
	if len(os.Getenv("INTEGRATION")) > 0 {
		t.Log("SKIPPING UNIT TEST")
		return
	}

	// TODO: add tests when we read the req body
	type args struct {
		flagKey string
	}

	type want struct {
		httpCode int
		bodyFile string
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "bare-minimum-flag",
			args: args{
				flagKey: "bare-minimum-flag",
			},
			want: want{
				httpCode: http.StatusOK,
				bodyFile: "../testdata/flag_eval/bare_minimum_response.json",
			},
		},
		{
			name: "bare-minimum-flag-2",
			args: args{
				flagKey: "bare-minimum-flag-2",
			},
			want: want{
				httpCode: http.StatusOK,
				bodyFile: "../testdata/flag_eval/bare_minimum_response_2.json",
			},
		},
		{
			name: "get default value if flag is disabled",
			args: args{
				flagKey: "disabled-flag",
			},
			want: want{
				httpCode: http.StatusOK,
				bodyFile: "../testdata/flag_eval/disabled_flag_response.json",
			},
		},
		{
			name: "valid flag",
			args: args{
				flagKey: "allow-access",
			},
			want: want{
				httpCode: http.StatusOK,
				bodyFile: "../testdata/flag_eval/valid_response.json",
			},
		},
		{
			name: "request non-existant flag",
			args: args{
				flagKey: "not-found",
			},
			want: want{
				httpCode: http.StatusNotFound,
				bodyFile: "../testdata/flag_eval/flag_not_found_response.json",
			},
		},
	}

	for _, test := range tests {
		// env vars
		os.Setenv("API_KEYS", tok)
		os.Setenv("READ_CLIENT_LOCATION", dir+"/flags.json")

		// config
		config.New()

		// resource
		name := config.Name()

		// log
		logBuffer := memoryutils.NewBuffer()

		logger := memorylog.NewLog(
			log.LogWithPrefix(name),
			memorylog.LogWithBuffer(logBuffer),
		)

		log.SetLogger(logger)

		// traces

		// metrics

		// clients
		writeClient := noop.NewWriter(
			writer.WithLocation(config.WriteClientLocation()),
		)

		readClient := localreader.NewReader(
			reader.WithLocation(config.ReadClientLocation()),
		)

		exportClient := localexporter.NewExporter(
			exporter.WithDir(config.ExportClientDir()),
		)

		notifyClient := localnotifier.NewNotifier()

		// servers
		httpServer, _, exportService, notifyService, err := server.Factory(
			writeClient,
			readClient,
			exportClient,
			notifyClient,
		)
		require.NoError(t, err)

		t.Run(test.name, func(t *testing.T) {
			err = httpServer.Run()
			require.NoError(t, err)

			req, err := http.NewRequest(
				http.MethodPost,
				fmt.Sprintf("http://%s%s%s", httpServer.Options().Address, "/ofrep/v1/evaluate/flags/", test.args.flagKey),
				strings.NewReader(""),
			)
			require.NoError(t, err)

			req.Header.Set("content-type", "application/json")
			req.Header.Set("authorization", fmt.Sprintf("Bearer %s", tok))

			client := &http.Client{}

			rsp, err := client.Do(req)
			require.NoError(t, err)

			want, err := os.ReadFile(test.want.bodyFile)
			require.NoError(t, err)

			got, err := io.ReadAll(rsp.Body)
			require.NoError(t, err)

			require.Equal(t, string(want), string(got))

			require.Equal(t, test.want.httpCode, rsp.StatusCode)

			t.Cleanup(func() {
				rsp.Body.Close()
				exportService.Close()
				notifyService.Close()
				config.Reset()
			})
		})
	}
}
