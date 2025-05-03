package file

import (
	"encoding/json"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

type Parser struct {
}

func (p *Parser) ParseFlags(bs []byte) (map[string]Flag, error) {
	var flags map[string]Flag
	var err error

	// TODO: from config
	switch strings.ToLower("yaml") {
	case "json":
		err = json.Unmarshal(bs, &flags)
	default:
		err = yaml.Unmarshal(bs, &flags)
	}

	for _, flag := range flags {
		if len(flag.Variations) == 0 {
			return nil, fmt.Errorf("flag missing variations")
		}

		if flag.DefaultRule == nil {
			return nil, fmt.Errorf("flag missing default rule")
		}

		if err := p.ParseRule(*flag.DefaultRule, true); err != nil {
			return nil, err
		}

		ruleNames := map[string]any{}

		for _, rule := range flag.Rules {
			if err := p.ParseRule(rule, false); err != nil {
				return nil, err
			}

			if _, ok := ruleNames[rule.Name]; ok {
				return nil, fmt.Errorf("multiple rules with the same name")
			} else {
				ruleNames[rule.Name] = nil
			}
		}

		flagValue, _ := flag.Evaluate(EvaluationContext{
			DefaultValue: nil,
		})

		switch v := flagValue; v.(type) {
		case int, float64, bool, string:
		default:
			if v != nil {
				return nil, fmt.Errorf("flag value %+v is not supported", v)
			}
		}
	}

	return flags, err
}

func (p *Parser) ParseRule(rule Rule, _ bool) error {
	if len(rule.Variation) == 0 {
		return fmt.Errorf("rule missing variation")
	}

	return nil
}
