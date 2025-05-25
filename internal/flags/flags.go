package flags

import (
	"encoding/json"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

func Factory(bs []byte, format string) (map[string]*Flag, error) {
	flags := map[string]*Flag{}

	var err error

	switch strings.ToLower(format) {
	case "json":
		err = json.Unmarshal(bs, &flags)
	default:
		err = yaml.Unmarshal(bs, &flags)
	}

	if err != nil {
		return nil, err
	}

	for k, flag := range flags {
		// sanity checks
		if flag == nil {
			return nil, fmt.Errorf("nil flag")
		}

		if len(k) == 0 {
			return nil, fmt.Errorf("flag missing key")
		}

		// add the default
		flag.DefaultRule = &Rule{
			Name:    "default",
			Variant: "default",
		}

		// requirements
		if len(flag.Variants) == 0 {
			return nil, fmt.Errorf("flag missing variants")
		}

		if _, ok := flag.Variants["default"]; !ok {
			return nil, fmt.Errorf("flag missing default variant")
		}

		ruleNames := map[string]any{}

		for _, rule := range flag.Rules {
			if rule == nil {
				return nil, fmt.Errorf("nil rule")
			}

			if err := parseRule(rule, flag.Variants); err != nil {
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

		for _, variant := range flag.Variants {
			var currentType string

			if len(variantType) > 0 {
				currentType, err = extractVariantType(variant)
				if err != nil {
					return nil, err
				}
				if currentType != variantType {
					return nil, fmt.Errorf("discovered flag variants with different types")
				}
			} else {
				variantType, err = extractVariantType(variant)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	return flags, nil
}

func parseRule(rule *Rule, variants map[string]any) error {
	if len(rule.Name) == 0 {
		return fmt.Errorf("rule missing name")
	}

	// we need to have exactly one of these but not both
	// 1) variant (with or without query)
	// 2) percentages (with variants that add up to 100)

	if len(rule.Variant) == 0 {
		return fmt.Errorf("rule missing variant")
	}

	// if this thing has percentages, check the variants there instead
	if _, ok := variants[rule.Variant]; !ok {
		return fmt.Errorf("rule includes unknown variant")
	}

	return nil
}

func extractVariantType(variant any) (string, error) {
	switch variant.(type) {
	case int, float64:
		return "number", nil
	case bool:
		return "bool", nil
	case string:
		return "string", nil
	default:
		return "", fmt.Errorf("flag value %+v is not supported", variant)
	}
}
