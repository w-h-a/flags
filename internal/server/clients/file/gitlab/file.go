package gitlab

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/w-h-a/flags/internal/server/clients/file"
)

type client struct {
	options    file.Options
	httpClient *http.Client
	parser     *file.Parser
}

func (c *client) Read(ctx context.Context) (map[string]*file.Flag, error) {
	// TODO: generalize
	file := c.options.Files[0]

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("https://gitlab.com/api/v4/projects/%s/repository/files/%s/raw?ref=main", c.options.Dir, file),
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

	return c.parser.ParseFlags(bs)
}

func NewFileClient(opts ...file.Option) file.Client {
	options := file.NewOptions(opts...)

	if err := options.Validate(); err != nil {
		log.Fatalf("failed to configure gitlab client: %v", err)
	}

	httpClient := http.DefaultClient
	httpClient.Timeout = 10 * time.Second

	c := &client{
		options:    options,
		httpClient: httpClient,
		parser:     &file.Parser{},
	}

	return c
}
