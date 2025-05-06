package allflags

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

func TestAllFlags(t *testing.T) {
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
				bodyFile: "../testdata/all_flags/valid_response.json",
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
				fmt.Sprintf("http://%s%s", httpServer.Options().Address, "/ofrep/v1/evaluate/flags"),
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
