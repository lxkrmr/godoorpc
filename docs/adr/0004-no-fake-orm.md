# ADR 0004: No fake ORM

## Status

Accepted

## Context

The OCA `odoorpc` Python library wraps Odoo's JSON-RPC API behind an
interface that mimics the Odoo ORM: `odoo.env["res.partner"]`,
`model.browse(ids)`, record attribute access. This is familiar to Odoo
developers but creates a misleading mental model.

The caller is not inside Odoo. There is no ORM. Every attribute access,
every `browse`, every field read is a remote procedure call over HTTP.
The abstraction hides that fact and makes it harder to reason about what
is actually happening on the wire.

`godoorpc` is a Go library targeting developers who work with Odoo
tooling, not Odoo application developers writing business logic. The
relevant mental model is HTTP and RPC, not Django-style ORM.

## Decision

`godoorpc` does not provide ORM-style abstractions. There are no model
proxies, no `browse`, no record objects, no attribute access that
triggers a network call.

The public API is `ExecuteKW(model, method, args, kwargs)`. Callers know
they are making a remote call. The result is a plain Go value.

## Consequences

Positive:
- honest mental model: every call is an explicit RPC call
- no surprise network activity hidden behind attribute access
- simpler library with less surface area to maintain

Negative:
- more verbose than a fluent ORM-style API
- callers must know Odoo model names and method signatures directly
