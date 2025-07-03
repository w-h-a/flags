package parseflags

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/w-h-a/flags/internal/flags"
)

func TestParseFlags(t *testing.T) {
	if len(os.Getenv("INTEGRATION")) > 0 {
		t.Log("SKIPPING UNIT TEST")
		return
	}

	tests := []struct {
		name     string
		filePath string
		format   string
		wantErr  bool
		err      string
	}{
		{
			name:     "valid yaml",
			filePath: "../testdata/flags.yaml",
			format:   "yaml",
			wantErr:  false,
		},
		{
			name:     "valid json ",
			filePath: "../testdata/flags.json",
			format:   "json",
			wantErr:  false,
		},
		{
			name:     "valid yaml upper",
			filePath: "../testdata/flags.yaml",
			format:   "YAML",
			wantErr:  false,
		},
		{
			name:     "valid json upper",
			filePath: "../testdata/flags.json",
			format:   "JSON",
			wantErr:  false,
		},
		{
			name:     "no variants yaml",
			filePath: "../testdata/parse_flags/no_variants.yaml",
			format:   "yaml",
			wantErr:  true,
			err:      "flag missing variants",
		},
		{
			name:     "no variants json",
			filePath: "../testdata/parse_flags/no_variants.json",
			format:   "json",
			wantErr:  true,
			err:      "flag missing variants",
		},
		{
			name:     "no default yaml",
			filePath: "../testdata/parse_flags/no_default.yaml",
			format:   "yaml",
			wantErr:  true,
			err:      "flag missing default variant",
		},
		{
			name:     "no default json",
			filePath: "../testdata/parse_flags/no_default.json",
			format:   "json",
			wantErr:  true,
			err:      "flag missing default variant",
		},
		{
			name:     "different variant types yaml",
			filePath: "../testdata/parse_flags/different_variant_types.yaml",
			format:   "yaml",
			wantErr:  true,
			err:      "discovered flag variants with different types",
		},
		{
			name:     "different variant types json",
			filePath: "../testdata/parse_flags/different_variant_types.json",
			format:   "json",
			wantErr:  true,
			err:      "discovered flag variants with different types",
		},
		{
			name:     "no name rule yaml",
			filePath: "../testdata/parse_flags/no_name_rule.yaml",
			format:   "yaml",
			wantErr:  true,
			err:      "rule missing name",
		},
		{
			name:     "no name rule json",
			filePath: "../testdata/parse_flags/no_name_rule.json",
			format:   "json",
			wantErr:  true,
			err:      "rule missing name",
		},
		{
			name:     "same name rules yaml",
			filePath: "../testdata/parse_flags/same_name_rules.yaml",
			format:   "yaml",
			wantErr:  true,
			err:      "multiple rules with the same name",
		},
		{
			name:     "same name rules json",
			filePath: "../testdata/parse_flags/same_name_rules.json",
			format:   "json",
			wantErr:  true,
			err:      "multiple rules with the same name",
		},
		{
			name:     "unknown variant rule yaml",
			filePath: "../testdata/parse_flags/unknown_variant_rule.yaml",
			format:   "yaml",
			wantErr:  true,
			err:      "rule includes unknown variant",
		},
		{
			name:     "unknown variant rule json",
			filePath: "../testdata/parse_flags/unknown_variant_rule.json",
			format:   "json",
			wantErr:  true,
			err:      "rule includes unknown variant",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			bs, err := os.ReadFile(test.filePath)
			require.NoError(t, err)

			_, err = flags.Factory(bs, test.format)

			if test.wantErr {
				require.Error(t, err)
				require.Equal(t, test.err, err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
