package report

import "context"

type Client interface {
	Create(ctx context.Context, events []Record) error
}
