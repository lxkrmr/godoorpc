package godoorpc

import (
	"encoding/json"
	"fmt"
	"strings"
)

// ParseDomain parses an Odoo domain string into a Domain.
//
// The input uses Python-style syntax with tuples and prefix operators:
//
//	[('is_company', '=', True)]
//	['|', ('name', 'ilike', 'foo'), ('name', 'ilike', 'bar')]
//
// Supported Python literals: True, False, None, strings, integers, floats.
func ParseDomain(s string) (Domain, error) {
	s = toJSON(s)

	var raw []json.RawMessage
	if err := json.Unmarshal([]byte(s), &raw); err != nil {
		return nil, fmt.Errorf("parse domain %q: %w", s, err)
	}

	domain := make(Domain, 0, len(raw))
	for _, item := range raw {
		node, err := parseDomainNode(item)
		if err != nil {
			return nil, fmt.Errorf("parse domain node: %w", err)
		}
		domain = append(domain, node)
	}
	return domain, nil
}

// toJSON converts Python-style domain syntax to valid JSON.
func toJSON(s string) string {
	s = strings.ReplaceAll(s, "(", "[")
	s = strings.ReplaceAll(s, ")", "]")
	s = strings.ReplaceAll(s, "'", "\"")
	s = strings.ReplaceAll(s, "True", "true")
	s = strings.ReplaceAll(s, "False", "false")
	s = strings.ReplaceAll(s, "None", "null")
	return s
}

// parseDomainNode turns a raw JSON value into a DomainNode.
// A JSON string becomes an Operator, a JSON array becomes a Condition.
func parseDomainNode(raw json.RawMessage) (DomainNode, error) {
	// Try string first — could be an operator like "|", "&", "!"
	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		return Operator(s), nil
	}

	// Otherwise expect a three-element array: [field, operator, value]
	var triple []json.RawMessage
	if err := json.Unmarshal(raw, &triple); err != nil {
		return nil, fmt.Errorf("expected a string operator or a [field, op, value] triple, got: %s", raw)
	}
	if len(triple) != 3 {
		return nil, fmt.Errorf("condition must have exactly 3 elements, got %d: %s", len(triple), raw)
	}

	var field, operator string
	if err := json.Unmarshal(triple[0], &field); err != nil {
		return nil, fmt.Errorf("condition field must be a string: %w", err)
	}
	if err := json.Unmarshal(triple[1], &operator); err != nil {
		return nil, fmt.Errorf("condition operator must be a string: %w", err)
	}

	var value any
	if err := json.Unmarshal(triple[2], &value); err != nil {
		return nil, fmt.Errorf("condition value: %w", err)
	}

	return Condition{Field: field, Operator: operator, Value: value}, nil
}
