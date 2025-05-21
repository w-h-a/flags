package exporter

import "context"

type Exporter interface {
	Export(ctx context.Context, events []Record) error
}
