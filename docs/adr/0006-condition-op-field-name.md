# ADR 0006: Condition.Op instead of Condition.Operator

## Status

Accepted

## Context

`Condition` holds a comparison triple: field, comparison operator, value.
The package also defines an `Operator` type for logical prefix operators
(`|`, `&`, `!`) used in domain expressions.

Naming the field `Condition.Operator string` creates confusion: why is
this field a plain `string` while there is also a type called `Operator`?
The answer — they represent different concepts — is not visible from the
field name alone.

## Decision

The field is named `Op`, not `Operator`.

`Op` is a common Go abbreviation for operator. It clearly refers to a
comparison operator (`=`, `!=`, `>`, `ilike`, etc.) without suggesting
any relation to the `Operator` type, which is reserved for logical prefix
operators in domain expressions.

## Consequences

Positive:
- no naming clash between `Condition.Op` and the `Operator` type
- the two concepts — comparison operator and logical prefix operator —
  are visually distinct at every call site

Negative:
- `Op` is shorter than `Operator` and may feel abbreviated to some readers
