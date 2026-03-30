# ADR 0003: Zero external dependencies

## Status

Accepted

## Context

`godoorpc` communicates with Odoo over HTTP using JSON. Go's standard
library provides `net/http` for HTTP and `encoding/json` for JSON. No
third-party packages are needed for the core functionality.

External dependencies carry ongoing costs: supply chain risk, version
drift, transitive dependencies, and maintenance burden. The Python
ecosystem has seen repeated supply chain attacks via PyPI. Keeping the
dependency surface at zero eliminates that class of risk entirely.

Tools built on `godoorpc` also benefit: a library with no dependencies
does not pull unexpected packages into the tool's own module graph.

## Decision

`godoorpc` uses only the Go standard library. No external packages are
added to `go.mod` beyond the Go toolchain itself.

If a future need genuinely cannot be met by the standard library, the
decision is revisited with a concrete justification.

## Consequences

Positive:
- no supply chain risk from third-party packages
- no transitive dependencies
- `go.sum` stays empty
- simple, auditable codebase

Negative:
- some tasks require more code than a library would need
- contributors cannot reach for well-known packages as shortcuts
