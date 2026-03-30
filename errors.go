package godoorpc

import "fmt"

// AuthError is returned when Odoo rejects the login credentials.
type AuthError struct {
	Message string
}

func (e *AuthError) Error() string {
	return e.Message
}

// RPCError is returned when Odoo responds with a JSON-RPC error.
type RPCError struct {
	Code    int
	Message string
}

func (e *RPCError) Error() string {
	return fmt.Sprintf("rpc error %d: %s", e.Code, e.Message)
}
