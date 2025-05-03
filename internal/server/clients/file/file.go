package file

import (
	"context"
	"errors"
)

var (
	ErrRuleDoesNotApply = errors.New("rule does not apply")
)

type Client interface {
	Read(ctx context.Context) (map[string]*Flag, error)
}
