# ADR 0007: Strict validation of domain prefix operators

## Status

Accepted

## Context

In `ParseDomain`, domain nodes that are JSON strings become prefix
operators. Without validation, any string — including typos like `"||"`
or `"or"` — would silently produce an `Operator` value and reach Odoo
as part of the RPC call, causing a confusing server-side error or
silently wrong behavior.

Odoo defines exactly three prefix operators for domain expressions:
`|` (OR), `&` (AND), and `!` (NOT). This set is stable and unlikely
to change.

## Decision

`ParseDomain` validates that any string node is one of the three known
operators. An unknown string returns an error immediately with a clear
message naming the invalid value and the valid alternatives.

## Consequences

Positive:
- typos in operator strings produce a clear error at parse time
- callers get a specific error message instead of a silent wrong domain

Negative:
- if a future Odoo version introduces a new prefix operator, ParseDomain
  will reject it until the library is updated
