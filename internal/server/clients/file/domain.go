package file

const (
	ReasonDisabled       string = "DISABLED"
	ReasonDefault        string = "DEFAULT"
	ReasonTargetingMatch string = "TARGETING_MATCH"

	ErrorNotFound string = "FLAG_NOT_FOUND"
)

type Flag struct {
	Disabled *bool           `json:"disabled" yaml:"disabled"`
	Variants map[string]*any `json:"variants" yaml:"variants"`
	Rules    []*Rule         `json:"rules" yaml:"rules"`

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
				Reason:    ReasonTargetingMatch,
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
		return true
	}

	return *f.Disabled
}

func (f *Flag) value(name string) any {
	for k, v := range f.Variants {
		if k == name && v != nil {
			return *v
		}
	}

	return nil
}

type Rule struct {
	Name    string `json:"name" yaml:"name"`
	Variant string `json:"variant" yaml:"variant"`
}

func (r *Rule) Evaluate(_ bool) string {
	return r.Variant
}

type ResolutionDetails struct {
	Variant   string
	Reason    string
	RuleIndex int
	RuleName  string
}
