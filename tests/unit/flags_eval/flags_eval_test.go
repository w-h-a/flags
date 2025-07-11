package flagseval

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
)

const (
	tok = "mytoken"
	dir = "../testdata"
)

func TestAllFlags_YAML(t *testing.T) {
	if len(os.Getenv("INTEGRATION")) > 0 {
		t.Log("SKIPPING UNIT TEST")
		return
	}

	// TODO: add tests when we read the req body
	type want struct {
		httpCode int
		bodyFile string
	}

	tests := []struct {
		name string
		want want
	}{
		{
			name: "valid flags",
			want: want{
				httpCode: http.StatusOK,
				bodyFile: "../testdata/flags_eval/valid_response.json",
			},
		},
	}

	for _, test := range tests {
		// env vars
		os.Setenv("API_KEYS", tok)
		os.Setenv("READ_CLIENT_LOCATION", dir+"/flags.yaml")

		// config
		config.New()

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
				fmt.Sprintf("http://%s%s", httpServer.Options().Address, "/ofrep/v1/evaluate/flags"),
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
				notifyService.Close()
				exportService.Close()
				err = httpServer.Stop()
				require.NoError(t, err)
				config.Reset()
			})
		})
	}
}

func TestAllFlags_JSON(t *testing.T) {
	if len(os.Getenv("INTEGRATION")) > 0 {
		t.Log("SKIPPING UNIT TEST")
		return
	}

	// TODO: add tests when we read the req body
	type want struct {
		httpCode int
		bodyFile string
	}

	tests := []struct {
		name string
		want want
	}{
		{
			name: "valid flags",
			want: want{
				httpCode: http.StatusOK,
				bodyFile: "../testdata/flags_eval/valid_response.json",
			},
		},
	}

	for _, test := range tests {
		// env vars
		os.Setenv("API_KEYS", tok)
		os.Setenv("READ_CLIENT_LOCATION", dir+"/flags.json")

		// config
		config.New()

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
				fmt.Sprintf("http://%s%s", httpServer.Options().Address, "/ofrep/v1/evaluate/flags"),
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
				notifyService.Close()
				exportService.Close()
				err = httpServer.Stop()
				require.NoError(t, err)
				config.Reset()
			})
		})
	}
}
