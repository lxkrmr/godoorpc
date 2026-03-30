package godoorpc

import "encoding/json"

// DomainNode is either a Condition or an Operator.
type DomainNode interface {
	isDomainNode()
}

// Condition represents a single filter triple: (field, op, value).
// Op is a comparison operator such as "=", "!=", ">", "<", "ilike".
// This is distinct from Operator, which is a logical prefix operator.
type Condition struct {
	Field string
	Op    string
	Value any
}

func (c Condition) isDomainNode() {}

// MarshalJSON serializes Condition as a JSON array [field, op, value].
func (c Condition) MarshalJSON() ([]byte, error) {
	return json.Marshal([]any{c.Field, c.Op, c.Value})
}

// Operator is a prefix domain operator.
type Operator string

func (o Operator) isDomainNode() {}

const (
	Or  Operator = "|"
	And Operator = "&"
	Not Operator = "!"
)

// Domain is a list of DomainNodes used as a filter in Odoo RPC calls.
type Domain []DomainNode

// MarshalJSON serializes Domain as a JSON array.
// A nil or empty Domain marshals to [] — never null.
// Odoo expects an empty array for an unrestricted domain, not null.
func (d Domain) MarshalJSON() ([]byte, error) {
	if len(d) == 0 {
		return []byte("[]"), nil
	}
	nodes := make([]any, len(d))
	for i, node := range d {
		nodes[i] = node
	}
	return json.Marshal(nodes)
}

// Args is a positional argument list for an RPC call.
type Args []any

// KWArgs is a keyword argument map for an RPC call.
type KWArgs map[string]any
