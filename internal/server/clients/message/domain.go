package message

import "github.com/w-h-a/flags/internal/server/clients/file"

type Diff struct {
	Deleted map[string]*file.Flag  `json:"deleted"`
	Added   map[string]*file.Flag  `json:"added"`
	Updated map[string]DiffUpdated `json:"updated"`
}

func (d *Diff) HasDiff() bool {
	return len(d.Deleted) > 0 || len(d.Added) > 0 || len(d.Updated) > 0
}

type DiffUpdated struct {
	Before *file.Flag `json:"old_value"`
	After  *file.Flag `json:"new_value"`
}
