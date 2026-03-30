# godoorpc

A minimal Go client library for the Odoo JSON-RPC API.

## Overview

`godoorpc` provides a session-based connection to Odoo and a single RPC
primitive that covers all operations: `ExecuteKW`. It stays close to the
wire — no fake ORM, no proxy objects, no hidden network calls.

## Install

```sh
go get github.com/lxkrmr/godoorpc
```

## Usage

```go
import "github.com/lxkrmr/godoorpc"

// Connect and authenticate
session, err := godoorpc.NewSession("http://localhost:8069", "mydb", "admin", "password")
if err != nil {
    // handle error
}

// Call any Odoo model method
result, err := session.ExecuteKW("res.partner", "search_read",
    godoorpc.Args{
        godoorpc.Domain{
            godoorpc.Condition{Field: "is_company", Operator: "=", Value: true},
        },
    },
    godoorpc.KWArgs{
        "fields": []string{"name", "email"},
        "limit":  10,
    },
)

// Parse a domain from a string (useful for CLI tools)
domain, err := godoorpc.ParseDomain("[('is_company', '=', True)]")
```

## Design

- **One RPC primitive** — everything goes through `ExecuteKW`, which maps
  directly to Odoo's `/web/dataset/call_kw` endpoint.
- **Zero external dependencies** — only Go stdlib.
- **Typed domain** — `Domain`, `Condition`, and `Operator` model Odoo's
  domain syntax with proper Go types.

See `docs/adr/` for the reasoning behind key design decisions.

## License

MIT
