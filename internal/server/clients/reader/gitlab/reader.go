package gitlab

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/w-h-a/flags/internal/flags"
	"github.com/w-h-a/flags/internal/server/clients/reader"
	"github.com/w-h-a/flags/internal/server/config"
	"github.com/w-h-a/pkg/telemetry/log"
)

type client struct {
	options    reader.Options
	httpClient *http.Client
}

func (c *client) Read(ctx context.Context) (map[string]*flags.Flag, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("https://gitlab.com/api/v4/projects/%s/repository/files/%s/raw?ref=main", c.options.Dir, c.options.File),
		strings.NewReader(""),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	if len(c.options.Token) > 0 {
		req.Header.Add("private-token", c.options.Token)
	}

	rsp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}

	defer rsp.Body.Close()

	if rsp.StatusCode > 399 {
		return nil, fmt.Errorf("received status code %d from github", rsp.StatusCode)
	}

	bs, err := io.ReadAll(rsp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body from github: %v", err)
	}

	return flags.Factory(bs, config.FlagFormat())
}

func NewReader(opts ...reader.Option) reader.Reader {
	options := reader.NewOptions(opts...)

	if err := options.Validate(); err != nil {
		log.Fatalf("failed to configure gitlab client: %v", err)
	}

	httpClient := http.DefaultClient
	httpClient.Timeout = 10 * time.Second

	c := &client{
		options:    options,
		httpClient: httpClient,
	}

	return c
}
