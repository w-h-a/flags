package go_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/open-feature/go-sdk-contrib/providers/ofrep"
	of "github.com/open-feature/go-sdk/openfeature"
	"github.com/stretchr/testify/require"
)

const (
	tok = "mytoken"
)

func TestGo_BooleanEval(t *testing.T) {
	if len(os.Getenv("INTEGRATION")) == 0 || len(os.Getenv("DYNAMODB")) > 0 {
		t.Log("SKIPPING INTEGRATION TEST")
		return
	}

	type args struct {
		apiKey       string
		flag         string
		defaultValue bool
		evalCtx      of.EvaluationContext
	}

	tests := []struct {
		name string
		args args
		want of.BooleanEvaluationDetails
	}{
		{
			name: "resolve a boolean flag with TARGETING_MATCH reason",
			args: args{
				apiKey:       tok,
				flag:         "bool_targeting_match",
				defaultValue: false,
				evalCtx: of.NewEvaluationContext(
					"",
					map[string]any{},
				),
			},
			want: of.BooleanEvaluationDetails{
				Value: true,
				EvaluationDetails: of.EvaluationDetails{
					FlagKey:  "bool_targeting_match",
					FlagType: of.Boolean,
					ResolutionDetail: of.ResolutionDetail{
						Variant:      "true",
						Reason:       of.TargetingMatchReason,
						ErrorCode:    "",
						ErrorMessage: "",
						FlagMetadata: map[string]any{},
					},
				},
			},
		},
		{
			name: "resolve with DEFAULT reason",
			args: args{
				apiKey:       tok,
				flag:         "default_bool",
				defaultValue: true,
				evalCtx: of.NewEvaluationContext(
					"",
					map[string]any{},
				),
			},
			want: of.BooleanEvaluationDetails{
				Value: false,
				EvaluationDetails: of.EvaluationDetails{
					FlagKey:  "default_bool",
					FlagType: of.Boolean,
					ResolutionDetail: of.ResolutionDetail{
						Variant:      "default",
						Reason:       of.DefaultReason,
						ErrorCode:    "",
						ErrorMessage: "",
						FlagMetadata: map[string]any{},
					},
				},
			},
		},
		{
			name: "use default if the flag is disabled",
			args: args{
				apiKey:       tok,
				flag:         "disabled_bool",
				defaultValue: false,
				evalCtx: of.NewEvaluationContext(
					"",
					map[string]any{},
				),
			},
			want: of.BooleanEvaluationDetails{
				Value: false,
				EvaluationDetails: of.EvaluationDetails{
					FlagKey:  "disabled_bool",
					FlagType: of.Boolean,
					ResolutionDetail: of.ResolutionDetail{
						Variant:      "default",
						Reason:       of.DisabledReason,
						ErrorCode:    "",
						ErrorMessage: "",
						FlagMetadata: map[string]any{},
					},
				},
			},
		},
		{
			name: "use default if the flag is disabled 2",
			args: args{
				apiKey:       tok,
				flag:         "disabled_bool_2",
				defaultValue: false,
				evalCtx: of.NewEvaluationContext(
					"",
					map[string]any{},
				),
			},
			want: of.BooleanEvaluationDetails{
				Value: false,
				EvaluationDetails: of.EvaluationDetails{
					FlagKey:  "disabled_bool_2",
					FlagType: of.Boolean,
					ResolutionDetail: of.ResolutionDetail{
						Variant:      "default",
						Reason:       of.DisabledReason,
						ErrorCode:    "",
						ErrorMessage: "",
						FlagMetadata: map[string]any{},
					},
				},
			},
		},
		{
			name: "error if we expect a boolean but get another type",
			args: args{
				apiKey:       tok,
				flag:         "string_targeting_match",
				defaultValue: false,
				evalCtx: of.NewEvaluationContext(
					"",
					map[string]any{},
				),
			},
			want: of.BooleanEvaluationDetails{
				Value: false,
				EvaluationDetails: of.EvaluationDetails{
					FlagKey:  "string_targeting_match",
					FlagType: of.Boolean,
					ResolutionDetail: of.ResolutionDetail{
						Variant:      "",
						Reason:       of.ErrorReason,
						ErrorCode:    of.TypeMismatchCode,
						ErrorMessage: "resolved value fdsa is not of boolean type",
						FlagMetadata: map[string]any{},
					},
				},
			},
		},
		{
			name: "error if there is no flag",
			args: args{
				apiKey:       tok,
				flag:         "no_such_flag",
				defaultValue: false,
				evalCtx: of.NewEvaluationContext(
					"",
					map[string]any{},
				),
			},
			want: of.BooleanEvaluationDetails{
				Value: false,
				EvaluationDetails: of.EvaluationDetails{
					FlagKey:  "no_such_flag",
					FlagType: of.Boolean,
					ResolutionDetail: of.ResolutionDetail{
						Variant:      "",
						Reason:       of.ErrorReason,
						ErrorCode:    of.FlagNotFoundCode,
						ErrorMessage: "flag for key 'no_such_flag' does not exist",
						FlagMetadata: map[string]any{},
					},
				},
			},
		},
		{
			name: "resolve to default with invalid api key",
			args: args{
				apiKey:       "notthedroid",
				flag:         "bool_targeting_match",
				defaultValue: false,
				evalCtx: of.NewEvaluationContext(
					"",
					map[string]any{},
				),
			},
			want: of.BooleanEvaluationDetails{
				Value: false,
				EvaluationDetails: of.EvaluationDetails{
					FlagKey:  "bool_targeting_match",
					FlagType: of.Boolean,
					ResolutionDetail: of.ResolutionDetail{
						Variant:      "",
						Reason:       of.ErrorReason,
						ErrorCode:    of.GeneralCode,
						ErrorMessage: "authentication/authorization error",
						FlagMetadata: map[string]any{},
					},
				},
			},
		},
		{
			name: "resolve to default with no api key",
			args: args{
				apiKey:       "",
				flag:         "bool_targeting_match",
				defaultValue: false,
				evalCtx: of.NewEvaluationContext(
					"",
					map[string]any{},
				),
			},
			want: of.BooleanEvaluationDetails{
				Value: false,
				EvaluationDetails: of.EvaluationDetails{
					FlagKey:  "bool_targeting_match",
					FlagType: of.Boolean,
					ResolutionDetail: of.ResolutionDetail{
						Variant:      "",
						Reason:       of.ErrorReason,
						ErrorCode:    of.GeneralCode,
						ErrorMessage: "authentication/authorization error",
						FlagMetadata: map[string]any{},
					},
				},
			},
		},
		{
			name: "resolve to true when targetingKey does match",
			args: args{
				apiKey:       tok,
				flag:         "bool_query",
				defaultValue: false,
				evalCtx: of.NewEvaluationContext(
					"123456",
					map[string]any{},
				),
			},
			want: of.BooleanEvaluationDetails{
				Value: true,
				EvaluationDetails: of.EvaluationDetails{
					FlagKey:  "bool_query",
					FlagType: of.Boolean,
					ResolutionDetail: of.ResolutionDetail{
						Variant:      "true",
						Reason:       of.TargetingMatchReason,
						ErrorCode:    "",
						ErrorMessage: "",
						FlagMetadata: map[string]any{},
					},
				},
			},
		},
		{
			name: "resolve to false when targetingKey does NOT match but there is a backup rule",
			args: args{
				apiKey:       tok,
				flag:         "bool_query",
				defaultValue: false,
				evalCtx: of.NewEvaluationContext(
					"654321",
					map[string]any{},
				),
			},
			want: of.BooleanEvaluationDetails{
				Value: false,
				EvaluationDetails: of.EvaluationDetails{
					FlagKey:  "bool_query",
					FlagType: of.Boolean,
					ResolutionDetail: of.ResolutionDetail{
						Variant:      "false",
						Reason:       of.TargetingMatchReason,
						ErrorCode:    "",
						ErrorMessage: "",
						FlagMetadata: map[string]any{},
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			provider := ofrep.NewProvider(
				fmt.Sprintf("http://%s:%d", "localhost", 4000),
				ofrep.WithBearerToken(test.args.apiKey),
			)

			err := of.SetProviderAndWait(provider)
			require.NoError(t, err)

			client := of.NewClient("test")

			got, err := client.BooleanValueDetails(
				context.TODO(),
				test.args.flag,
				test.args.defaultValue,
				test.args.evalCtx,
			)

			if len(test.want.ErrorCode) > 0 {
				require.Error(t, err)
				require.Equal(t, fmt.Sprintf("error code: %s: %s", test.want.ErrorCode, test.want.ErrorMessage), err.Error())
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, test.want, got)
		})
	}
}

func TestGo_FloatEval(t *testing.T) {
	if len(os.Getenv("INTEGRATION")) == 0 || len(os.Getenv("DYNAMODB")) > 0 {
		t.Log("SKIPPING INTEGRATION TEST")
		return
	}

	type args struct {
		apiKey       string
		flag         string
		defaultValue float64
		evalCtx      of.EvaluationContext
	}

	tests := []struct {
		name string
		args args
		want of.FloatEvaluationDetails
	}{
		{
			name: "resolve a float flag with TARGETING_MATCH reason",
			args: args{
				apiKey:       tok,
				flag:         "float_targeting_match",
				defaultValue: 0.0,
				evalCtx: of.NewEvaluationContext(
					"",
					map[string]any{},
				),
			},
			want: of.FloatEvaluationDetails{
				Value: 101.25,
				EvaluationDetails: of.EvaluationDetails{
					FlagKey:  "float_targeting_match",
					FlagType: of.Float,
					ResolutionDetail: of.ResolutionDetail{
						Variant:      "true",
						Reason:       of.TargetingMatchReason,
						ErrorCode:    "",
						ErrorMessage: "",
						FlagMetadata: map[string]any{},
					},
				},
			},
		},
		{
			name: "resolve with DEFAULT reason",
			args: args{
				apiKey:       tok,
				flag:         "default_float",
				defaultValue: 0.0,
				evalCtx: of.NewEvaluationContext(
					"",
					map[string]any{},
				),
			},
			want: of.FloatEvaluationDetails{
				Value: 100.25,
				EvaluationDetails: of.EvaluationDetails{
					FlagKey:  "default_float",
					FlagType: of.Float,
					ResolutionDetail: of.ResolutionDetail{
						Variant:      "default",
						Reason:       of.DefaultReason,
						ErrorCode:    "",
						ErrorMessage: "",
						FlagMetadata: map[string]any{},
					},
				},
			},
		},
		{
			name: "use default if the flag is disabled",
			args: args{
				apiKey:       tok,
				flag:         "disabled_float",
				defaultValue: 100.25,
				evalCtx: of.NewEvaluationContext(
					"",
					map[string]any{},
				),
			},
			want: of.FloatEvaluationDetails{
				Value: 100.25,
				EvaluationDetails: of.EvaluationDetails{
					FlagKey:  "disabled_float",
					FlagType: of.Float,
					ResolutionDetail: of.ResolutionDetail{
						Variant:      "default",
						Reason:       of.DisabledReason,
						ErrorCode:    "",
						ErrorMessage: "",
						FlagMetadata: map[string]any{},
					},
				},
			},
		},
		{
			name: "use default if the flag is disabled 2",
			args: args{
				apiKey:       tok,
				flag:         "disabled_float_2",
				defaultValue: 100.25,
				evalCtx: of.NewEvaluationContext(
					"",
					map[string]any{},
				),
			},
			want: of.FloatEvaluationDetails{
				Value: 100.25,
				EvaluationDetails: of.EvaluationDetails{
					FlagKey:  "disabled_float_2",
					FlagType: of.Float,
					ResolutionDetail: of.ResolutionDetail{
						Variant:      "default",
						Reason:       of.DisabledReason,
						ErrorCode:    "",
						ErrorMessage: "",
						FlagMetadata: map[string]any{},
					},
				},
			},
		},
		{
			name: "error if we expect a float but get another type",
			args: args{
				apiKey:       tok,
				flag:         "string_targeting_match",
				defaultValue: 0.0,
				evalCtx: of.NewEvaluationContext(
					"",
					map[string]any{},
				),
			},
			want: of.FloatEvaluationDetails{
				Value: 0.0,
				EvaluationDetails: of.EvaluationDetails{
					FlagKey:  "string_targeting_match",
					FlagType: of.Float,
					ResolutionDetail: of.ResolutionDetail{
						Variant:      "",
						Reason:       of.ErrorReason,
						ErrorCode:    of.TypeMismatchCode,
						ErrorMessage: "resolved value fdsa is not of float type",
						FlagMetadata: map[string]any{},
					},
				},
			},
		},
		{
			name: "error if there is no flag",
			args: args{
				apiKey:       tok,
				flag:         "no_such_flag",
				defaultValue: 0.0,
				evalCtx: of.NewEvaluationContext(
					"",
					map[string]any{},
				),
			},
			want: of.FloatEvaluationDetails{
				Value: 0.0,
				EvaluationDetails: of.EvaluationDetails{
					FlagKey:  "no_such_flag",
					FlagType: of.Float,
					ResolutionDetail: of.ResolutionDetail{
						Variant:      "",
						Reason:       of.ErrorReason,
						ErrorCode:    of.FlagNotFoundCode,
						ErrorMessage: "flag for key 'no_such_flag' does not exist",
						FlagMetadata: map[string]any{},
					},
				},
			},
		},
		{
			name: "resolve to default with invalid api key",
			args: args{
				apiKey:       "notthedroid",
				flag:         "float_targeting_match",
				defaultValue: 0.0,
				evalCtx: of.NewEvaluationContext(
					"",
					map[string]any{},
				),
			},
			want: of.FloatEvaluationDetails{
				Value: 0.0,
				EvaluationDetails: of.EvaluationDetails{
					FlagKey:  "float_targeting_match",
					FlagType: of.Float,
					ResolutionDetail: of.ResolutionDetail{
						Variant:      "",
						Reason:       of.ErrorReason,
						ErrorCode:    of.GeneralCode,
						ErrorMessage: "authentication/authorization error",
						FlagMetadata: map[string]any{},
					},
				},
			},
		},
		{
			name: "resolve to default with no api key",
			args: args{
				apiKey:       "",
				flag:         "float_targeting_match",
				defaultValue: 0.0,
				evalCtx: of.NewEvaluationContext(
					"",
					map[string]any{},
				),
			},
			want: of.FloatEvaluationDetails{
				Value: 0.0,
				EvaluationDetails: of.EvaluationDetails{
					FlagKey:  "float_targeting_match",
					FlagType: of.Float,
					ResolutionDetail: of.ResolutionDetail{
						Variant:      "",
						Reason:       of.ErrorReason,
						ErrorCode:    of.GeneralCode,
						ErrorMessage: "authentication/authorization error",
						FlagMetadata: map[string]any{},
					},
				},
			},
		},
		{
			name: "resolve to 100.10 when targetingKey does match",
			args: args{
				apiKey:       tok,
				flag:         "float_query",
				defaultValue: 0.0,
				evalCtx: of.NewEvaluationContext(
					"123456",
					map[string]any{},
				),
			},
			want: of.FloatEvaluationDetails{
				Value: 100.1,
				EvaluationDetails: of.EvaluationDetails{
					FlagKey:  "float_query",
					FlagType: of.Float,
					ResolutionDetail: of.ResolutionDetail{
						Variant:      "true",
						Reason:       of.TargetingMatchReason,
						ErrorCode:    "",
						ErrorMessage: "",
						FlagMetadata: map[string]any{},
					},
				},
			},
		},
		{
			name: "resolve to 100.00 when targetingKey does NOT match but there is a backup rule",
			args: args{
				apiKey:       tok,
				flag:         "float_query",
				defaultValue: 0.0,
				evalCtx: of.NewEvaluationContext(
					"654321",
					map[string]any{},
				),
			},
			want: of.FloatEvaluationDetails{
				Value: 100.00,
				EvaluationDetails: of.EvaluationDetails{
					FlagKey:  "float_query",
					FlagType: of.Float,
					ResolutionDetail: of.ResolutionDetail{
						Variant:      "false",
						Reason:       of.TargetingMatchReason,
						ErrorCode:    "",
						ErrorMessage: "",
						FlagMetadata: map[string]any{},
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			provider := ofrep.NewProvider(
				fmt.Sprintf("http://%s:%d", "localhost", 4000),
				ofrep.WithBearerToken(test.args.apiKey),
			)

			err := of.SetProviderAndWait(provider)
			require.NoError(t, err)

			client := of.NewClient("test")

			got, err := client.FloatValueDetails(
				context.TODO(),
				test.args.flag,
				test.args.defaultValue,
				test.args.evalCtx,
			)

			if len(test.want.ErrorCode) > 0 {
				require.Error(t, err)
				require.Equal(t, fmt.Sprintf("error code: %s: %s", test.want.ErrorCode, test.want.ErrorMessage), err.Error())
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, test.want, got)
		})
	}
}

func TestGo_IntEval(t *testing.T) {
	if len(os.Getenv("INTEGRATION")) == 0 || len(os.Getenv("DYNAMODB")) > 0 {
		t.Log("SKIPPING INTEGRATION TEST")
		return
	}

	type args struct {
		apiKey       string
		flag         string
		defaultValue int64
		evalCtx      of.EvaluationContext
	}

	tests := []struct {
		name string
		args args
		want of.IntEvaluationDetails
	}{
		{
			name: "resolve a integer flag with TARGETING_MATCH reason",
			args: args{
				apiKey:       tok,
				flag:         "int_targeting_match",
				defaultValue: 0,
				evalCtx: of.NewEvaluationContext(
					"",
					map[string]any{},
				),
			},
			want: of.IntEvaluationDetails{
				Value: 101,
				EvaluationDetails: of.EvaluationDetails{
					FlagKey:  "int_targeting_match",
					FlagType: of.Int,
					ResolutionDetail: of.ResolutionDetail{
						Variant:      "true",
						Reason:       of.TargetingMatchReason,
						ErrorCode:    "",
						ErrorMessage: "",
						FlagMetadata: map[string]any{},
					},
				},
			},
		},
		{
			name: "resolve with DEFAULT reason",
			args: args{
				apiKey:       tok,
				flag:         "default_int",
				defaultValue: 0,
				evalCtx: of.NewEvaluationContext(
					"",
					map[string]any{},
				),
			},
			want: of.IntEvaluationDetails{
				Value: 100,
				EvaluationDetails: of.EvaluationDetails{
					FlagKey:  "default_int",
					FlagType: of.Int,
					ResolutionDetail: of.ResolutionDetail{
						Variant:      "default",
						Reason:       of.DefaultReason,
						ErrorCode:    "",
						ErrorMessage: "",
						FlagMetadata: map[string]any{},
					},
				},
			},
		},
		{
			name: "use default if the flag is disabled",
			args: args{
				apiKey:       tok,
				flag:         "disabled_int",
				defaultValue: 100,
				evalCtx: of.NewEvaluationContext(
					"",
					map[string]any{},
				),
			},
			want: of.IntEvaluationDetails{
				Value: 100,
				EvaluationDetails: of.EvaluationDetails{
					FlagKey:  "disabled_int",
					FlagType: of.Int,
					ResolutionDetail: of.ResolutionDetail{
						Variant:      "default",
						Reason:       of.DisabledReason,
						ErrorCode:    "",
						ErrorMessage: "",
						FlagMetadata: map[string]any{},
					},
				},
			},
		},
		{
			name: "use default if the flag is disabled 2",
			args: args{
				apiKey:       tok,
				flag:         "disabled_int_2",
				defaultValue: 100,
				evalCtx: of.NewEvaluationContext(
					"",
					map[string]any{},
				),
			},
			want: of.IntEvaluationDetails{
				Value: 100,
				EvaluationDetails: of.EvaluationDetails{
					FlagKey:  "disabled_int_2",
					FlagType: of.Int,
					ResolutionDetail: of.ResolutionDetail{
						Variant:      "default",
						Reason:       of.DisabledReason,
						ErrorCode:    "",
						ErrorMessage: "",
						FlagMetadata: map[string]any{},
					},
				},
			},
		},
		{
			name: "error if we expect a int but get another type",
			args: args{
				apiKey:       tok,
				flag:         "string_targeting_match",
				defaultValue: 0,
				evalCtx: of.NewEvaluationContext(
					"",
					map[string]any{},
				),
			},
			want: of.IntEvaluationDetails{
				Value: 0,
				EvaluationDetails: of.EvaluationDetails{
					FlagKey:  "string_targeting_match",
					FlagType: of.Int,
					ResolutionDetail: of.ResolutionDetail{
						Variant:      "",
						Reason:       of.ErrorReason,
						ErrorCode:    of.TypeMismatchCode,
						ErrorMessage: "resolved value fdsa is not of integer type",
						FlagMetadata: map[string]any{},
					},
				},
			},
		},
		{
			name: "error if there is no flag",
			args: args{
				apiKey:       tok,
				flag:         "no_such_flag",
				defaultValue: 0,
				evalCtx: of.NewEvaluationContext(
					"",
					map[string]any{},
				),
			},
			want: of.IntEvaluationDetails{
				Value: 0,
				EvaluationDetails: of.EvaluationDetails{
					FlagKey:  "no_such_flag",
					FlagType: of.Int,
					ResolutionDetail: of.ResolutionDetail{
						Variant:      "",
						Reason:       of.ErrorReason,
						ErrorCode:    of.FlagNotFoundCode,
						ErrorMessage: "flag for key 'no_such_flag' does not exist",
						FlagMetadata: map[string]any{},
					},
				},
			},
		},
		{
			name: "resolve to default with invalid api key",
			args: args{
				apiKey:       "notthedroid",
				flag:         "int_targeting_match",
				defaultValue: 0,
				evalCtx: of.NewEvaluationContext(
					"",
					map[string]any{},
				),
			},
			want: of.IntEvaluationDetails{
				Value: 0,
				EvaluationDetails: of.EvaluationDetails{
					FlagKey:  "int_targeting_match",
					FlagType: of.Int,
					ResolutionDetail: of.ResolutionDetail{
						Variant:      "",
						Reason:       of.ErrorReason,
						ErrorCode:    of.GeneralCode,
						ErrorMessage: "authentication/authorization error",
						FlagMetadata: map[string]any{},
					},
				},
			},
		},
		{
			name: "resolve to default with no api key",
			args: args{
				apiKey:       "",
				flag:         "int_targeting_match",
				defaultValue: 0.0,
				evalCtx: of.NewEvaluationContext(
					"",
					map[string]any{},
				),
			},
			want: of.IntEvaluationDetails{
				Value: 0,
				EvaluationDetails: of.EvaluationDetails{
					FlagKey:  "int_targeting_match",
					FlagType: of.Int,
					ResolutionDetail: of.ResolutionDetail{
						Variant:      "",
						Reason:       of.ErrorReason,
						ErrorCode:    of.GeneralCode,
						ErrorMessage: "authentication/authorization error",
						FlagMetadata: map[string]any{},
					},
				},
			},
		},
		{
			name: "resolve to 101 when targetingKey does match",
			args: args{
				apiKey:       tok,
				flag:         "int_query",
				defaultValue: 0,
				evalCtx: of.NewEvaluationContext(
					"123456",
					map[string]any{},
				),
			},
			want: of.IntEvaluationDetails{
				Value: 101,
				EvaluationDetails: of.EvaluationDetails{
					FlagKey:  "int_query",
					FlagType: of.Int,
					ResolutionDetail: of.ResolutionDetail{
						Variant:      "true",
						Reason:       of.TargetingMatchReason,
						ErrorCode:    "",
						ErrorMessage: "",
						FlagMetadata: map[string]any{},
					},
				},
			},
		},
		{
			name: "resolve to 100 when targetingKey does NOT match but there is a backup rule",
			args: args{
				apiKey:       tok,
				flag:         "int_query",
				defaultValue: 0.0,
				evalCtx: of.NewEvaluationContext(
					"654321",
					map[string]any{},
				),
			},
			want: of.IntEvaluationDetails{
				Value: 100.00,
				EvaluationDetails: of.EvaluationDetails{
					FlagKey:  "int_query",
					FlagType: of.Int,
					ResolutionDetail: of.ResolutionDetail{
						Variant:      "false",
						Reason:       of.TargetingMatchReason,
						ErrorCode:    "",
						ErrorMessage: "",
						FlagMetadata: map[string]any{},
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			provider := ofrep.NewProvider(
				fmt.Sprintf("http://%s:%d", "localhost", 4000),
				ofrep.WithBearerToken(test.args.apiKey),
			)

			err := of.SetProviderAndWait(provider)
			require.NoError(t, err)

			client := of.NewClient("test")

			got, err := client.IntValueDetails(
				context.TODO(),
				test.args.flag,
				test.args.defaultValue,
				test.args.evalCtx,
			)

			if len(test.want.ErrorCode) > 0 {
				require.Error(t, err)
				require.Equal(t, fmt.Sprintf("error code: %s: %s", test.want.ErrorCode, test.want.ErrorMessage), err.Error())
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, test.want, got)
		})
	}
}

func TestGo_StringEval(t *testing.T) {
	if len(os.Getenv("INTEGRATION")) == 0 || len(os.Getenv("DYNAMODB")) > 0 {
		t.Log("SKIPPING INTEGRATION TEST")
		return
	}

	type args struct {
		apiKey       string
		flag         string
		defaultValue string
		evalCtx      of.EvaluationContext
	}

	tests := []struct {
		name string
		args args
		want of.StringEvaluationDetails
	}{
		{
			name: "resolve a string flag with TARGETING_MATCH reason",
			args: args{
				apiKey:       tok,
				flag:         "string_targeting_match",
				defaultValue: "",
				evalCtx: of.NewEvaluationContext(
					"",
					map[string]any{},
				),
			},
			want: of.StringEvaluationDetails{
				Value: "fdsa",
				EvaluationDetails: of.EvaluationDetails{
					FlagKey:  "string_targeting_match",
					FlagType: of.String,
					ResolutionDetail: of.ResolutionDetail{
						Variant:      "true",
						Reason:       of.TargetingMatchReason,
						ErrorCode:    "",
						ErrorMessage: "",
						FlagMetadata: map[string]any{},
					},
				},
			},
		},
		{
			name: "resolve with DEFAULT reason",
			args: args{
				apiKey:       tok,
				flag:         "default_string",
				defaultValue: "",
				evalCtx: of.NewEvaluationContext(
					"",
					map[string]any{},
				),
			},
			want: of.StringEvaluationDetails{
				Value: "asdf",
				EvaluationDetails: of.EvaluationDetails{
					FlagKey:  "default_string",
					FlagType: of.String,
					ResolutionDetail: of.ResolutionDetail{
						Variant:      "default",
						Reason:       of.DefaultReason,
						ErrorCode:    "",
						ErrorMessage: "",
						FlagMetadata: map[string]any{},
					},
				},
			},
		},
		{
			name: "use default if the flag is disabled",
			args: args{
				apiKey:       tok,
				flag:         "disabled_string",
				defaultValue: "asdf",
				evalCtx: of.NewEvaluationContext(
					"",
					map[string]any{},
				),
			},
			want: of.StringEvaluationDetails{
				Value: "asdf",
				EvaluationDetails: of.EvaluationDetails{
					FlagKey:  "disabled_string",
					FlagType: of.String,
					ResolutionDetail: of.ResolutionDetail{
						Variant:      "default",
						Reason:       of.DisabledReason,
						ErrorCode:    "",
						ErrorMessage: "",
						FlagMetadata: map[string]any{},
					},
				},
			},
		},
		{
			name: "use default if the flag is disabled 2",
			args: args{
				apiKey:       tok,
				flag:         "disabled_string_2",
				defaultValue: "asdf",
				evalCtx: of.NewEvaluationContext(
					"",
					map[string]any{},
				),
			},
			want: of.StringEvaluationDetails{
				Value: "asdf",
				EvaluationDetails: of.EvaluationDetails{
					FlagKey:  "disabled_string_2",
					FlagType: of.String,
					ResolutionDetail: of.ResolutionDetail{
						Variant:      "default",
						Reason:       of.DisabledReason,
						ErrorCode:    "",
						ErrorMessage: "",
						FlagMetadata: map[string]any{},
					},
				},
			},
		},
		{
			name: "error if we expect a string but get another type",
			args: args{
				apiKey:       tok,
				flag:         "int_targeting_match",
				defaultValue: "",
				evalCtx: of.NewEvaluationContext(
					"",
					map[string]any{},
				),
			},
			want: of.StringEvaluationDetails{
				Value: "",
				EvaluationDetails: of.EvaluationDetails{
					FlagKey:  "int_targeting_match",
					FlagType: of.String,
					ResolutionDetail: of.ResolutionDetail{
						Variant:      "",
						Reason:       of.ErrorReason,
						ErrorCode:    of.TypeMismatchCode,
						ErrorMessage: "resolved value 101 is not of string type",
						FlagMetadata: map[string]any{},
					},
				},
			},
		},
		{
			name: "error if there is no flag",
			args: args{
				apiKey:       tok,
				flag:         "no_such_flag",
				defaultValue: "",
				evalCtx: of.NewEvaluationContext(
					"",
					map[string]any{},
				),
			},
			want: of.StringEvaluationDetails{
				Value: "",
				EvaluationDetails: of.EvaluationDetails{
					FlagKey:  "no_such_flag",
					FlagType: of.String,
					ResolutionDetail: of.ResolutionDetail{
						Variant:      "",
						Reason:       of.ErrorReason,
						ErrorCode:    of.FlagNotFoundCode,
						ErrorMessage: "flag for key 'no_such_flag' does not exist",
						FlagMetadata: map[string]any{},
					},
				},
			},
		},
		{
			name: "resolve to default with invalid api key",
			args: args{
				apiKey:       "notthedroid",
				flag:         "string_targeting_match",
				defaultValue: "",
				evalCtx: of.NewEvaluationContext(
					"",
					map[string]any{},
				),
			},
			want: of.StringEvaluationDetails{
				Value: "",
				EvaluationDetails: of.EvaluationDetails{
					FlagKey:  "string_targeting_match",
					FlagType: of.String,
					ResolutionDetail: of.ResolutionDetail{
						Variant:      "",
						Reason:       of.ErrorReason,
						ErrorCode:    of.GeneralCode,
						ErrorMessage: "authentication/authorization error",
						FlagMetadata: map[string]any{},
					},
				},
			},
		},
		{
			name: "resolve to default with no api key",
			args: args{
				apiKey:       "",
				flag:         "string_targeting_match",
				defaultValue: "",
				evalCtx: of.NewEvaluationContext(
					"",
					map[string]any{},
				),
			},
			want: of.StringEvaluationDetails{
				Value: "",
				EvaluationDetails: of.EvaluationDetails{
					FlagKey:  "string_targeting_match",
					FlagType: of.String,
					ResolutionDetail: of.ResolutionDetail{
						Variant:      "",
						Reason:       of.ErrorReason,
						ErrorCode:    of.GeneralCode,
						ErrorMessage: "authentication/authorization error",
						FlagMetadata: map[string]any{},
					},
				},
			},
		},
		{
			name: "resolve to 'fdsa' when targetingKey does match",
			args: args{
				apiKey:       tok,
				flag:         "string_query",
				defaultValue: "",
				evalCtx: of.NewEvaluationContext(
					"123456",
					map[string]any{},
				),
			},
			want: of.StringEvaluationDetails{
				Value: "fdsa",
				EvaluationDetails: of.EvaluationDetails{
					FlagKey:  "string_query",
					FlagType: of.String,
					ResolutionDetail: of.ResolutionDetail{
						Variant:      "true",
						Reason:       of.TargetingMatchReason,
						ErrorCode:    "",
						ErrorMessage: "",
						FlagMetadata: map[string]any{},
					},
				},
			},
		},
		{
			name: "resolve to 'asdf' when targetingKey does NOT match but there is a backup rule",
			args: args{
				apiKey:       tok,
				flag:         "string_query",
				defaultValue: "",
				evalCtx: of.NewEvaluationContext(
					"654321",
					map[string]any{},
				),
			},
			want: of.StringEvaluationDetails{
				Value: "asdf",
				EvaluationDetails: of.EvaluationDetails{
					FlagKey:  "string_query",
					FlagType: of.String,
					ResolutionDetail: of.ResolutionDetail{
						Variant:      "false",
						Reason:       of.TargetingMatchReason,
						ErrorCode:    "",
						ErrorMessage: "",
						FlagMetadata: map[string]any{},
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			provider := ofrep.NewProvider(
				fmt.Sprintf("http://%s:%d", "localhost", 4000),
				ofrep.WithBearerToken(test.args.apiKey),
			)

			err := of.SetProviderAndWait(provider)
			require.NoError(t, err)

			client := of.NewClient("test")

			got, err := client.StringValueDetails(
				context.TODO(),
				test.args.flag,
				test.args.defaultValue,
				test.args.evalCtx,
			)

			if len(test.want.ErrorCode) > 0 {
				require.Error(t, err)
				require.Equal(t, fmt.Sprintf("error code: %s: %s", test.want.ErrorCode, test.want.ErrorMessage), err.Error())
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, test.want, got)
		})
	}
}
