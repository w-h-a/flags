package unit

import (
	"github.com/w-h-a/flags/internal/flags"
)

func DefaultFlags() map[string]*flags.Flag {
	return map[string]*flags.Flag{
		"flag1": {
			Disabled: Bool(false),
			Variants: map[string]any{
				"default":  "A",
				"variant2": "B",
			},
			Rules: []*flags.Rule{
				{
					Name:    "rule1",
					Variant: "variant2",
				},
			},
		},
		"flag2": {
			Disabled: Bool(false),
			Variants: map[string]any{
				"default":  "A",
				"variant2": "B",
			},
			Rules: []*flags.Rule{
				{
					Name:    "rule1",
					Variant: "variant2",
				},
			},
		},
	}
}

func Bool(v bool) *bool {
	return &v
}

func Float64(v float64) *float64 {
	return &v
}

func Int(v int) *int {
	return &v
}

func String(v string) *string {
	return &v
}
