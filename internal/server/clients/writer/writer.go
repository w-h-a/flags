package writer

import (
	"context"
)

type Writer interface {
	Write(ctx context.Context, key string, bs []byte) error
}
