package message

import (
	"context"
	"sync"
)

type Client interface {
	Send(ctx context.Context, diff Diff, waitGroup *sync.WaitGroup) error
}
