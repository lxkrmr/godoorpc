package godoorpc_test

import (
	"encoding/json"
	"testing"

	"github.com/lxkrmr/godoorpc"
)

func TestDomain_MarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		domain   godoorpc.Domain
		expected string
	}{
		{
			name:     "nil domain marshals to empty array",
			domain:   nil,
			expected: `[]`,
		},
		{
			name:     "empty domain marshals to empty array",
			domain:   godoorpc.Domain{},
			expected: `[]`,
		},
		{
			name: "single condition",
			domain: godoorpc.Domain{
				godoorpc.Condition{Field: "is_company", Op: "=", Value: true},
			},
			expected: `[["is_company","=",true]]`,
		},
		{
			name: "OR with two conditions",
			domain: godoorpc.Domain{
				godoorpc.Or,
				godoorpc.Condition{Field: "name", Op: "ilike", Value: "foo"},
				godoorpc.Condition{Field: "name", Op: "ilike", Value: "bar"},
			},
			expected: `["|",["name","ilike","foo"],["name","ilike","bar"]]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := json.Marshal(tt.domain)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if string(b) != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, string(b))
			}
		})
	}
}
