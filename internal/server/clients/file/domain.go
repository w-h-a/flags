package file

const (
	ReasonDisabled  string = "DISABLED"
	ReasonRuleMatch string = "RULE_MATCH"
	ReasonDefault   string = "DEFAULT"
	ReasonError     string = "ERROR"
)

type Flag struct {
	Disabled   *bool           `json:"disabled" yaml:"disabled"`
	Variations map[string]*any `json:"variations" yaml:"variations"`
	Rules      []*Rule         `json:"rules" yaml:"rules"`

	DefaultRule *Rule
}

func (f *Flag) Evaluate() (any, ResolutionDetails) {
	if f.IsDisabled() {
		variant := f.DefaultRule.Evaluate(true)

		resolutionDetails := ResolutionDetails{Variant: variant, Reason: ReasonDisabled}

		return f.value(variant), resolutionDetails
	}

	if len(f.Rules) > 0 {
		for i, rule := range f.Rules {
			variant := rule.Evaluate(false)

			resolutionDetails := ResolutionDetails{
				Variant:   variant,
				Reason:    ReasonRuleMatch,
				RuleIndex: i,
				RuleName:  rule.Name,
			}

			return f.value(variant), resolutionDetails
		}
	}

	variant := f.DefaultRule.Evaluate(true)

	resolutionDetails := ResolutionDetails{Variant: variant, Reason: ReasonDefault}

	return f.value(variant), resolutionDetails
}

func (f *Flag) IsDisabled() bool {
	if f.Disabled == nil {
		return false
	}

	return *f.Disabled
}

func (f *Flag) value(name string) any {
	for k, v := range f.Variations {
		if k == name && v != nil {
			return *v
		}
	}

	return nil
}

type Rule struct {
	Name      string `json:"name" yaml:"name"`
	Variation string `json:"variation" yaml:"variation"`
}

func (r *Rule) Evaluate(_ bool) string {
	return r.Variation
}

type ResolutionDetails struct {
	Variant   string
	Reason    string
	RuleIndex int
	RuleName  string
}
