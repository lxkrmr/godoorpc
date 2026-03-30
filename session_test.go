package godoorpc_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lxkrmr/godoorpc"
)

func rpcResult(result any) []byte {
	resp := map[string]any{"jsonrpc": "2.0", "id": 1, "result": result}
	b, _ := json.Marshal(resp)
	return b
}

func rpcErr(code int, message string) []byte {
	resp := map[string]any{
		"jsonrpc": "2.0",
		"id":      1,
		"error":   map[string]any{"code": code, "message": message},
	}
	b, _ := json.Marshal(resp)
	return b
}

func TestNewSession_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(rpcResult(map[string]any{"uid": 1, "username": "admin"}))
	}))
	defer srv.Close()

	_, err := godoorpc.NewSession(srv.URL, "testdb", "admin", "admin")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestNewSession_WrongPassword(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(rpcResult(map[string]any{"uid": false}))
	}))
	defer srv.Close()

	_, err := godoorpc.NewSession(srv.URL, "testdb", "admin", "wrong")
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	var authErr *godoorpc.AuthError
	if !errors.As(err, &authErr) {
		t.Fatalf("expected AuthError, got %T: %v", err, err)
	}
}

func TestNewSession_RPCError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(rpcErr(100, "database not found"))
	}))
	defer srv.Close()

	_, err := godoorpc.NewSession(srv.URL, "baddb", "admin", "admin")
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	var rpcErr *godoorpc.RPCError
	if !errors.As(err, &rpcErr) {
		t.Fatalf("expected RPCError, got %T: %v", err, err)
	}
}

func TestExecuteKW(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/web/session/authenticate":
			w.Write(rpcResult(map[string]any{"uid": 1}))
		case "/web/dataset/call_kw":
			w.Write(rpcResult([]any{
				map[string]any{"id": 1, "name": "Acme"},
			}))
		}
	}))
	defer srv.Close()

	session, err := godoorpc.NewSession(srv.URL, "testdb", "admin", "admin")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result, err := session.ExecuteKW(
		"res.partner", "search_read",
		godoorpc.Args{
			godoorpc.Domain{
				godoorpc.Condition{Field: "is_company", Operator: "=", Value: true},
			},
		},
		godoorpc.KWArgs{"fields": []string{"name"}, "limit": 1},
	)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result == nil {
		t.Fatal("expected result, got nil")
	}
}
