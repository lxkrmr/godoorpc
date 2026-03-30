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
// Bare Python literals (True, False, None) are converted to their JSON
// equivalents only when they appear outside of string values.
func toJSON(s string) string {
	s = strings.ReplaceAll(s, "(", "[")
	s = strings.ReplaceAll(s, ")", "]")
	s = strings.ReplaceAll(s, "'", "\"")
	// Replace bare Python literals — must come after quote replacement
	// so that string values like 'True' are already quoted as "True"
	// and won't be matched by a bare-word replace.
	s = replaceBare(s, "True", "true")
	s = replaceBare(s, "False", "false")
	s = replaceBare(s, "None", "null")
	return s
}

// replaceBare replaces occurrences of old with new only when they are not
// inside a JSON string (i.e. not surrounded by double quotes).
func replaceBare(s, old, new string) string {
	var b strings.Builder
	inString := false
	i := 0
	for i < len(s) {
		if s[i] == '"' {
			inString = !inString
			b.WriteByte(s[i])
			i++
			continue
		}
		if !inString && strings.HasPrefix(s[i:], old) {
			b.WriteString(new)
			i += len(old)
			continue
		}
		b.WriteByte(s[i])
		i++
	}
	return b.String()
}

// knownOperators holds the set of valid Odoo prefix domain operators.
var knownOperators = map[string]Operator{
	string(Or):  Or,
	string(And): And,
	string(Not): Not,
}

// parseDomainNode turns a raw JSON value into a DomainNode.
// A JSON string becomes an Operator if it is one of |, &, or !.
// A JSON array becomes a Condition.
func parseDomainNode(raw json.RawMessage) (DomainNode, error) {
	// Try string first — must be a known prefix operator
	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		op, ok := knownOperators[s]
		if !ok {
			return nil, fmt.Errorf("unknown domain operator %q: expected one of %q, %q, %q", s, Or, And, Not)
		}
		return op, nil
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

	return Condition{Field: field, Op: operator, Value: value}, nil
}
