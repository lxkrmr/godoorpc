# ADR 0001: call_kw as the single RPC primitive

## Status

Accepted

## Context

Odoo exposes multiple JSON-RPC endpoints. The two most relevant for
a developer client are `/web/dataset/call_kw` and `/web/dataset/search_read`.

`call_kw` is the universal endpoint. It covers model method calls of any
kind: `search_read`, `fields_get`, `read`, `write`, `create`, `unlink`,
and button actions such as `button_immediate_upgrade`. Every operation
that a developer tool needs can be expressed as a `call_kw` request.

A dedicated `search_read` convenience method could reduce verbosity for
the most common read operation. But it would introduce a second mental
model for callers: when do I use `SearchRead` vs `ExecuteKW`?

## Decision

`godoorpc` exposes a single RPC primitive: `ExecuteKW`.

No separate `SearchRead`, `FieldsGet`, or other convenience wrappers are
provided in the initial library. Callers express all operations through
`ExecuteKW(model, method, args, kwargs)`.

## Consequences

Positive:
- one mental model for all RPC operations
- the library stays close to the wire
- no abstraction layer that could drift from Odoo's actual API
- callers understand exactly what they are sending

Negative:
- more verbose call sites for common operations like `search_read`
- callers must know Odoo method names directly
