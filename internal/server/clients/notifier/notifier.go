package notifier

import (
	"context"
)

type Notifier interface {
	Notify(ctx context.Context, diff Diff) error
}
