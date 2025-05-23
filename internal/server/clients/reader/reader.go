package reader

import (
	"context"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type Reader interface {
	ReadByKey(ctx context.Context, key string) ([]byte, error)
	Read(ctx context.Context) ([]byte, error)
}
