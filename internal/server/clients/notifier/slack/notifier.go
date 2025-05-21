package slack

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gdexlab/go-render/render"
	difflib "github.com/r3labs/diff/v3"
	"github.com/w-h-a/flags/internal/server/clients/notifier"
	"github.com/w-h-a/pkg/telemetry/log"
)

const (
	colorDeleted = "#FF0000"
	colorAdded   = "#008000"
	colorUpdated = "#FFA500"
)

type client struct {
	options    notifier.Options
	httpClient *http.Client
}

func (c *client) Notify(ctx context.Context, diff notifier.Diff) error {
	msg := c.convert(diff)

	bs, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %v", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.options.URL,
		bytes.NewReader(bs),
	)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Add("content-type", "application/json")

	rsp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}

	defer rsp.Body.Close()

	if rsp.StatusCode > 399 {
		return fmt.Errorf("received status code %d from slack", rsp.StatusCode)
	}

	return nil
}

func (c *client) convert(diff notifier.Diff) slackMessage {
	attachments := c.convertDeleted(diff)
	attachments = append(attachments, c.convertAdded(diff)...)
	attachments = append(attachments, c.convertUpdated(diff)...)

	result := slackMessage{
		Text:         "Changes detected in feature flags",
		Attachements: attachments,
	}

	return result
}

func (c *client) convertDeleted(diff notifier.Diff) []attachment {
	attachments := []attachment{}

	emoji := "😵"

	for k := range diff.Deleted {
		attachment := attachment{
			Title: fmt.Sprintf("%s Flag \"%s\" deleted", emoji, k),
			Color: colorDeleted,
		}

		attachments = append(attachments, attachment)
	}

	return attachments
}

func (c *client) convertAdded(diff notifier.Diff) []attachment {
	attachments := []attachment{}

	emoji := "🥳"

	for k := range diff.Added {
		attachment := attachment{
			Title: fmt.Sprintf("%s Flag \"%s\" added", emoji, k),
			Color: colorAdded,
		}

		attachments = append(attachments, attachment)
	}

	return attachments
}

func (c *client) convertUpdated(diff notifier.Diff) []attachment {
	attachments := []attachment{}

	emoji := "✏️"

	for k, v := range diff.Updated {
		attachment := attachment{
			Title:  fmt.Sprintf("%s Flag \"%s\" updated", emoji, k),
			Color:  colorUpdated,
			Fields: []field{},
		}

		changelog, _ := difflib.Diff(v.Before, v.After, difflib.AllowTypeMismatch(true))

		for _, change := range changelog {
			if change.Type != "update" {
				continue
			}

			value := fmt.Sprintf("%s => %s", render.Render(change.From), render.Render(change.To))

			attachment.Fields = append(
				attachment.Fields,
				field{Title: strings.Join(change.Path, "."), Value: value},
			)
		}

		attachments = append(attachments, attachment)
	}

	return attachments
}

func NewNotifier(opts ...notifier.Option) notifier.Notifier {
	options := notifier.NewOptions(opts...)

	if len(options.URL) == 0 {
		log.Fatal("slack message client requires URL")
	}

	httpClient := http.DefaultClient
	httpClient.Timeout = 10 * time.Second

	c := &client{
		options:    options,
		httpClient: httpClient,
	}

	return c
}
