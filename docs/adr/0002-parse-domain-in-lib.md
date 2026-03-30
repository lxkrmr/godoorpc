# ADR 0002: ParseDomain lives in the library

## Status

Accepted

## Context

Odoo domains use a Python-style syntax with tuples and prefix operators:

```
[('is_company', '=', True), '|', ('name', 'ilike', 'foo')]
```

CLI tools built on `godoorpc` receive this syntax as a raw string from
the user and need to convert it into typed Go values before passing it
to `ExecuteKW`.

This parsing logic is not trivial. If left to each tool, every tool that
accepts domain input will implement its own version. That leads to
duplication, inconsistency, and diverging edge case handling.

## Decision

`godoorpc` provides `ParseDomain(s string) (Domain, error)`.

The function converts the Python-style domain string into a typed
`Domain` value. It lives in the library because the `Domain` type is
defined there, and because every tool that calls Odoo with domain filters
will need it.

## Consequences

Positive:
- one implementation, shared by all tools
- consistent handling of edge cases (True/False/None, prefix operators)
- callers get a typed value, not a raw string

Negative:
- the library now contains string parsing logic unrelated to HTTP or RPC
- a tool that never accepts domain strings from users still pulls in the parser
