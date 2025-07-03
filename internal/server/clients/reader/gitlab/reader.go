package gitlab

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/w-h-a/flags/internal/server/clients/reader"
)

type client struct {
	options    reader.Options
	httpClient *http.Client
}

func (c *client) ReadByKey(ctx context.Context, key string) ([]byte, error) {
	return nil, nil
}

func (c *client) Read(ctx context.Context) ([]byte, error) {
	// Gitlab location:
	// https://gitlab.com/api/v4/projects/:id/repository/files/:filePath/raw?ref=main
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		c.options.Location,
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

	return io.ReadAll(rsp.Body)
}

func NewReader(opts ...reader.Option) reader.Reader {
	options := reader.NewOptions(opts...)

	if err := options.Validate(); err != nil {
		slog.ErrorContext(context.Background(), "failed to configure gitlab reader", "error", err)
		os.Exit(1)
	}

	httpClient := http.DefaultClient
	httpClient.Timeout = 10 * time.Second

	c := &client{
		options:    options,
		httpClient: httpClient,
	}

	return c
}
