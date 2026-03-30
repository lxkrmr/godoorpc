package godoorpc_test

import (
	"testing"

	"github.com/lxkrmr/godoorpc"
)

func TestParseDomain(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected godoorpc.Domain
	}{
		{
			name:  "empty domain",
			input: "[]",
			expected: godoorpc.Domain{},
		},
		{
			name:  "single condition with string value",
			input: "[('name', '=', 'Acme')]",
			expected: godoorpc.Domain{
				godoorpc.Condition{Field: "name", Operator: "=", Value: "Acme"},
			},
		},
		{
			name:  "single condition with bool true",
			input: "[('is_company', '=', True)]",
			expected: godoorpc.Domain{
				godoorpc.Condition{Field: "is_company", Operator: "=", Value: true},
			},
		},
		{
			name:  "single condition with bool false",
			input: "[('active', '=', False)]",
			expected: godoorpc.Domain{
				godoorpc.Condition{Field: "active", Operator: "=", Value: false},
			},
		},
		{
			name:  "single condition with None",
			input: "[('parent_id', '=', None)]",
			expected: godoorpc.Domain{
				godoorpc.Condition{Field: "parent_id", Operator: "=", Value: nil},
			},
		},
		{
			name:  "single condition with integer value",
			input: "[('bid_price', '>', 0)]",
			expected: godoorpc.Domain{
				godoorpc.Condition{Field: "bid_price", Operator: ">", Value: float64(0)},
			},
		},
		{
			name:  "string value containing True is preserved as string",
			input: "[('state', '=', 'True')]",
			expected: godoorpc.Domain{
				godoorpc.Condition{Field: "state", Operator: "=", Value: "True"},
			},
		},
		{
			name:  "string value containing False is preserved as string",
			input: "[('state', '=', 'False')]",
			expected: godoorpc.Domain{
				godoorpc.Condition{Field: "state", Operator: "=", Value: "False"},
			},
		},
		{
			name:  "OR with two conditions",
			input: "['|', ('name', 'ilike', 'foo'), ('name', 'ilike', 'bar')]",
			expected: godoorpc.Domain{
				godoorpc.Or,
				godoorpc.Condition{Field: "name", Operator: "ilike", Value: "foo"},
				godoorpc.Condition{Field: "name", Operator: "ilike", Value: "bar"},
			},
		},
		{
			name:  "AND with two conditions",
			input: "['&', ('is_company', '=', True), ('active', '=', True)]",
			expected: godoorpc.Domain{
				godoorpc.And,
				godoorpc.Condition{Field: "is_company", Operator: "=", Value: true},
				godoorpc.Condition{Field: "active", Operator: "=", Value: true},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := godoorpc.ParseDomain(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(got) != len(tt.expected) {
				t.Fatalf("expected %d nodes, got %d", len(tt.expected), len(got))
			}
			for i, node := range got {
				if node != tt.expected[i] {
					t.Errorf("node[%d]: expected %v, got %v", i, tt.expected[i], node)
				}
			}
		})
	}
}

func TestParseDomain_InvalidInput(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{name: "condition with wrong length", input: "[('name', '=')]"},
		{name: "condition field not a string", input: "[[1, '=', 'foo']]"},
		{name: "malformed json after conversion", input: "[unclosed"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := godoorpc.ParseDomain(tt.input)
			if err == nil {
				t.Fatal("expected an error, got nil")
			}
		})
	}
}
