package flags

import (
	"errors"

	"github.com/w-h-a/pkg/telemetry/log"
)

const (
	ReasonDisabled       string = "DISABLED"
	ReasonDefault        string = "DEFAULT"
	ReasonTargetingMatch string = "TARGETING_MATCH"
	ReasonSplit          string = "SPLIT"

	ErrorNotFound string = "FLAG_NOT_FOUND"
)

var (
	ErrRuleDoesNotApply = errors.New("rule does not apply")
)

type Flag struct {
	Disabled *bool          `json:"disabled" yaml:"disabled"`
	Variants map[string]any `json:"variants" yaml:"variants"`
	Rules    []*Rule        `json:"rules" yaml:"rules"`

	DefaultRule *Rule
}

func (f *Flag) Evaluate() (any, ResolutionDetails) {
	if f.IsDisabled() {
		variant, _ := f.DefaultRule.Evaluate()

		resolutionDetails := ResolutionDetails{Variant: variant, Reason: ReasonDisabled}

		return f.value(variant), resolutionDetails
	}

	for i, rule := range f.Rules {
		variant, err := rule.Evaluate()
		if err != nil && errors.Is(err, ErrRuleDoesNotApply) {
			continue
		} else if err != nil {
			log.Errorf("unexpected error during rule evaluation: %v", err)
			continue
		}

		resolutionDetails := ResolutionDetails{
			Variant:   variant,
			RuleIndex: i,
			RuleName:  rule.Name,
		}

		// reason is determined by nature of rule
		// if the rule has percentages, then SPLIT
		// otherwise, TARGETING_MATCH
		resolutionDetails.Reason = ReasonTargetingMatch

		return f.value(variant), resolutionDetails
	}

	variant, _ := f.DefaultRule.Evaluate()

	resolutionDetails := ResolutionDetails{Variant: variant, Reason: ReasonDefault}

	return f.value(variant), resolutionDetails
}

func (f *Flag) IsDisabled() bool {
	if f.Disabled == nil {
		return true
	}

	return *f.Disabled
}

func (f *Flag) value(name string) any {
	for k, v := range f.Variants {
		if k == name && v != nil {
			return v
		}
	}

	return nil
}

type Rule struct {
	Name    string `json:"name" yaml:"name"`
	Variant string `json:"variant" yaml:"variant"`
}

func (r *Rule) Evaluate() (string, error) {
	// if this rule has percentages, use them

	// if this rule has a query, check whether it applies
	// if it does, return variant
	// otherwise, return empty string and ErrRuleDoesNotApply

	// otherwise, return the variant
	return r.Variant, nil
}

type ResolutionDetails struct {
	Variant   string
	Reason    string
	RuleIndex int
	RuleName  string
}

type Diff struct {
	Deleted map[string]*Flag       `json:"deleted"`
	Added   map[string]*Flag       `json:"added"`
	Updated map[string]DiffUpdated `json:"updated"`
}

func (d *Diff) HasDiff() bool {
	return len(d.Deleted) > 0 || len(d.Added) > 0 || len(d.Updated) > 0
}

type DiffUpdated struct {
	Before *Flag `json:"old_value"`
	After  *Flag `json:"new_value"`
}
