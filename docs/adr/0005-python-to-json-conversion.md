# ADR 0005: Python-to-JSON conversion for domain strings

## Status

Accepted

## Context

Odoo domain strings use Python syntax: tuples, single-quoted strings, and
bare Python literals (True, False, None):

    [('is_company', '=', True), '|', ('name', 'ilike', 'foo')]

Go's standard library has no Python parser. The options for converting
this syntax to typed Go values are:

1. **String replacement** — convert `(` to `[`, `'` to `"`, `True` to
   `true`, etc., then parse as JSON. Simple, no dependencies.
2. **Full Python literal parser** — parse the Python AST subset used in
   Odoo domains. Correct, but significantly more code.
3. **Accept JSON syntax only** — require callers to pass valid JSON
   domain syntax instead of Python. Loses ergonomics for CLI users.

## Decision

`ParseDomain` uses string replacement to convert Python syntax to JSON,
then calls `json.Unmarshal`. The replacement is aware of string
boundaries: bare `True`, `False`, and `None` are replaced only when
they appear outside of quoted string values. This handles the most
common case where a string value happens to contain one of these words.

## Consequences

Positive:
- zero external dependencies
- covers all Odoo domain patterns used in practice
- CLI tools can accept the same domain syntax Odoo developers already know

Negative:
- not a full Python parser — edge cases involving escaped quotes inside
  string values are not handled
- callers who need exact Python literal semantics should pre-convert to
  JSON themselves and skip ParseDomain
