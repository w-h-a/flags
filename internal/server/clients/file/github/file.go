package github

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/w-h-a/flags/internal/flags"
	"github.com/w-h-a/flags/internal/server/clients/file"
	"github.com/w-h-a/flags/internal/server/config"
	"github.com/w-h-a/pkg/telemetry/log"
)

type client struct {
	options    file.Options
	httpClient *http.Client
}

func (c *client) Read(ctx context.Context) (map[string]*flags.Flag, error) {
	// TODO: generalize
	file := c.options.Files[0]

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("https://api.github.com/repos/%s/contents/%s?ref=main", c.options.Dir, file),
		strings.NewReader(""),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Add("accept", "application/vnd.github.raw")

	if len(c.options.Token) > 0 {
		req.Header.Add("authorization", fmt.Sprintf("Bearer %s", c.options.Token))
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

func NewFileClient(opts ...file.Option) file.Client {
	options := file.NewOptions(opts...)

	if err := options.Validate(); err != nil {
		log.Fatalf("failed to configure github client: %v", err)
	}

	httpClient := http.DefaultClient
	httpClient.Timeout = 10 * time.Second

	c := &client{
		options:    options,
		httpClient: httpClient,
	}

	return c
}
