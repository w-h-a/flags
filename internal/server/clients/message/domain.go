package message

import (
	"github.com/w-h-a/flags/internal/flags"
)

type Diff struct {
	Deleted map[string]*flags.Flag `json:"deleted"`
	Added   map[string]*flags.Flag `json:"added"`
	Updated map[string]DiffUpdated `json:"updated"`
}

func (d *Diff) HasDiff() bool {
	return len(d.Deleted) > 0 || len(d.Added) > 0 || len(d.Updated) > 0
}

type DiffUpdated struct {
	Before *flags.Flag `json:"old_value"`
	After  *flags.Flag `json:"new_value"`
}
