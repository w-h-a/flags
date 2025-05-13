package message

import (
	"context"
)

type Client interface {
	Send(ctx context.Context, diff Diff) error
}
