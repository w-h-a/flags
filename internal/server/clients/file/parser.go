package file

import (
	"encoding/json"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

type Parser struct {
}

func (p *Parser) ParseFlags(bs []byte) (map[string]*Flag, error) {
	var flags map[string]*Flag
	var err error

	// TODO: from config
	switch strings.ToLower("yaml") {
	case "json":
		err = json.Unmarshal(bs, &flags)
	default:
		err = yaml.Unmarshal(bs, &flags)
	}

	for _, flag := range flags {
		// corrective actions
		flag.DefaultRule = &Rule{
			Name:      "default",
			Variation: "default",
		}

		// requirements
		if len(flag.Variations) == 0 {
			return nil, fmt.Errorf("flag missing variations")
		}

		if _, ok := flag.Variations["default"]; !ok {
			return nil, fmt.Errorf("flag's variations missing default value")
		}

		ruleNames := map[string]any{}

		for _, rule := range flag.Rules {
			if err := p.ParseRule(rule); err != nil {
				return nil, err
			}

			if _, ok := ruleNames[rule.Name]; ok {
				return nil, fmt.Errorf("multiple rules with the same name")
			} else {
				ruleNames[rule.Name] = nil
			}
		}

		// more complicated requirement checks
		flagValue, _ := flag.Evaluate()

		switch v := flagValue; v.(type) {
		case int, float64, bool, string:
		default:
			return nil, fmt.Errorf("flag value %+v is not supported", v)
		}
	}

	return flags, err
}

func (p *Parser) ParseRule(rule *Rule) error {
	if len(rule.Name) == 0 {
		return fmt.Errorf("rule missing name")
	}

	if len(rule.Variation) == 0 {
		return fmt.Errorf("rule missing variation")
	}

	return nil
}
