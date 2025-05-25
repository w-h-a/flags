package getflag

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/w-h-a/flags/internal/flags"
	"github.com/w-h-a/flags/internal/server"
	"github.com/w-h-a/flags/internal/server/clients/exporter"
	localexporter "github.com/w-h-a/flags/internal/server/clients/exporter/local"
	localnotifier "github.com/w-h-a/flags/internal/server/clients/notifier/local"
	"github.com/w-h-a/flags/internal/server/clients/writereader"
	mockwritereader "github.com/w-h-a/flags/internal/server/clients/writereader/mock"
	"github.com/w-h-a/flags/internal/server/config"
	"github.com/w-h-a/flags/tests/unit"
	"github.com/w-h-a/pkg/telemetry/log"
	memorylog "github.com/w-h-a/pkg/telemetry/log/memory"
	"github.com/w-h-a/pkg/utils/memoryutils"
	"gopkg.in/yaml.v3"
)

const (
	tok = "mytoken"
)

func TestGetFlag(t *testing.T) {
	if len(os.Getenv("INTEGRATION")) > 0 {
		t.Log("SKIPPING UNIT TEST")
		return
	}

	type inputs struct {
		flags        map[string]*flags.Flag
		key          string
		unauthorized bool
		headers      map[string]string
	}

	type want struct {
		httpCode int
		bodyFile string
	}

	tests := []struct {
		name   string
		inputs inputs
		want   want
	}{
		{
			name: "200 if found",
			inputs: inputs{
				flags:   unit.DefaultFlags(),
				key:     "flag2",
				headers: map[string]string{},
			},
			want: want{
				httpCode: http.StatusOK,
				bodyFile: "../testdata/get_flag/valid_response_flag2.json",
			},
		},
		{
			name: "404 if not found",
			inputs: inputs{
				flags:   unit.DefaultFlags(),
				key:     "flag99",
				headers: map[string]string{},
			},
			want: want{
				httpCode: http.StatusNotFound,
				bodyFile: "../testdata/not_found.json",
			},
		},
		{
			name: "500 if error",
			inputs: inputs{
				flags: unit.DefaultFlags(),
				key:   "flag2",
				headers: map[string]string{
					"error_read_by_key": "failed to read by key",
				},
			},
			want: want{
				httpCode: http.StatusInternalServerError,
				bodyFile: "../testdata/read_by_key_error.json",
			},
		},
		{
			name: "403 if unauthorized",
			inputs: inputs{
				unauthorized: true,
			},
			want: want{
				httpCode: http.StatusUnauthorized,
				bodyFile: "../testdata/unauthorized.json",
			},
		},
	}

	for _, test := range tests {
		// env vars
		os.Setenv("API_KEYS", tok)
		os.Setenv("FLAG_FORMAT", "yaml")
		os.Setenv("WRITE_CLIENT_LOCATION", "any")

		// config
		config.New()

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
		writereadClient := mockwritereader.NewWriteReader(
			writereader.WithLocation(config.WriteClientLocation()),
		)

		for k, v := range test.inputs.flags {
			bs, err := yaml.Marshal(map[string]*flags.Flag{
				k: v,
			})
			require.NoError(t, err)

			err = writereadClient.Write(context.TODO(), k, bs)
			require.NoError(t, err)
		}

		exportClient := localexporter.NewExporter(
			exporter.WithDir(config.ExportClientDir()),
		)

		notifyClient := localnotifier.NewNotifier()

		// servers and services
		httpServer, _, exportService, notifyService, err := server.Factory(
			writereadClient,
			writereadClient,
			exportClient,
			notifyClient,
		)
		require.NoError(t, err)

		t.Run(test.name, func(t *testing.T) {
			err = httpServer.Run()
			require.NoError(t, err)

			req, err := http.NewRequest(
				http.MethodGet,
				fmt.Sprintf("http://%s%s%s", httpServer.Options().Address, "/admin/v1/flags", "/"+test.inputs.key),
				strings.NewReader(""),
			)
			require.NoError(t, err)

			req.Header.Set("content-type", "application/json")

			if !test.inputs.unauthorized {
				req.Header.Set("authorization", fmt.Sprintf("Bearer %s", tok))
			}

			for k, v := range test.inputs.headers {
				req.Header.Set(k, v)
			}

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
