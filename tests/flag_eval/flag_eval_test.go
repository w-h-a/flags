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
	"github.com/w-h-a/flags/internal/server/clients/file"
	localfile "github.com/w-h-a/flags/internal/server/clients/file/local"
	localmessage "github.com/w-h-a/flags/internal/server/clients/message/local"
	"github.com/w-h-a/flags/internal/server/config"
	"github.com/w-h-a/pkg/telemetry/log"
	memorylog "github.com/w-h-a/pkg/telemetry/log/memory"
	"github.com/w-h-a/pkg/utils/memoryutils"
)

const (
	dir   = "../testdata"
	files = "/flags.yaml"
)

func TestFlagEval(t *testing.T) {
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
				bodyFile: "../testdata/flag_eval/flag_not_found_response.txt",
			},
		},
	}

	for _, test := range tests {
		// env vars
		os.Setenv("FILE_CLIENT_DIR", dir)
		os.Setenv("FILE_CLIENT_FILES", files)

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
		fileClient := localfile.NewFileClient(
			file.WithDir(config.FileClientDir()),
			file.WithFiles(config.FileClientFiles()...),
		)

		messageClient := localmessage.NewMessageClient()

		// servers
		httpServer, _, notifyService, err := server.Factory(
			fileClient,
			messageClient,
		)
		require.NoError(t, err)

		t.Run(test.name, func(t *testing.T) {
			err = httpServer.Run()
			require.NoError(t, err)

			rsp, err := http.Post(
				fmt.Sprintf("http://%s%s%s", httpServer.Options().Address, "/ofrep/v1/evaluate/flags/", test.args.flagKey),
				"application/json",
				strings.NewReader(""),
			)
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
				config.Reset()
			})
		})
	}
}
