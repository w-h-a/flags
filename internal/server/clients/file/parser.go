package file

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/w-h-a/flags/internal/server/config"
	"gopkg.in/yaml.v3"
)

type Parser struct {
}

func (p *Parser) ParseFlags(bs []byte) (map[string]*Flag, error) {
	var flags map[string]*Flag
	var err error

	switch strings.ToLower(config.FlagFormat()) {
	case "json":
		err = json.Unmarshal(bs, &flags)
	default:
		err = yaml.Unmarshal(bs, &flags)
	}

	if err != nil {
		return nil, err
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
		var variantType string
		var err error

		for _, variant := range flag.Variations {
			var currentType string

			if len(variantType) > 0 {
				currentType, err = p.extractVariantType(variant)
				if err != nil {
					return nil, err
				}
				if currentType != variantType {
					return nil, fmt.Errorf("variants have different types")
				}
			} else {
				variantType, err = p.extractVariantType(variant)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	return flags, nil
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

func (p *Parser) extractVariantType(variant *any) (string, error) {
	v := (*variant)

	switch v.(type) {
	case int, float64:
		return "number", nil
	case bool:
		return "bool", nil
	case string:
		return "string", nil
	default:
		return "", fmt.Errorf("flag value %+v is not supported", v)
	}
}
