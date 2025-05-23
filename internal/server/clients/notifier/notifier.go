package notifier

import (
	"context"

	"github.com/w-h-a/flags/internal/flags"
)

type Notifier interface {
	Notify(ctx context.Context, diff flags.Diff) error
}
