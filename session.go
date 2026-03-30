package godoorpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/cookiejar"
)

// Session is the interface for making RPC calls to Odoo.
type Session interface {
	ExecuteKW(model, method string, args Args, kwargs KWArgs) (any, error)
}

// Client holds an authenticated session to Odoo.
type Client struct {
	baseURL    string
	db         string
	httpClient *http.Client
}

// NewSession connects to Odoo at url, authenticates with the given
// credentials, and returns a Client ready to make RPC calls.
func NewSession(url, db, user, password string) (*Client, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("create cookie jar: %w", err)
	}
	c := &Client{
		baseURL:    url,
		db:         db,
		httpClient: &http.Client{Jar: jar},
	}
	if err := c.login(user, password); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Client) login(user, password string) error {
	params := map[string]any{
		"db":       c.db,
		"login":    user,
		"password": password,
	}
	var result map[string]any
	if err := c.call("/web/session/authenticate", params, &result); err != nil {
		return err
	}
	uid, exists := result["uid"]
	if !exists || uid == false {
		return &AuthError{Message: "authentication failed: wrong login or password"}
	}
	return nil
}

// ExecuteKW calls method on model via /web/dataset/call_kw.
func (c *Client) ExecuteKW(model, method string, args Args, kwargs KWArgs) (any, error) {
	params := map[string]any{
		"model":  model,
		"method": method,
		"args":   args,
		"kwargs": kwargs,
	}
	var result any
	if err := c.call("/web/dataset/call_kw", params, &result); err != nil {
		return nil, err
	}
	return result, nil
}

type rpcEnvelope struct {
	JSONRPC string `json:"jsonrpc"`
	Method  string `json:"method"`
	ID      int    `json:"id"`
	Params  any    `json:"params"`
}

type rpcResponse struct {
	Result json.RawMessage `json:"result"`
	Error  *rpcErrorBody   `json:"error"`
}

type rpcErrorBody struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (c *Client) call(endpoint string, params any, out any) error {
	body, err := json.Marshal(rpcEnvelope{
		JSONRPC: "2.0",
		Method:  "call",
		ID:      1,
		Params:  params,
	})
	if err != nil {
		return fmt.Errorf("marshal rpc request: %w", err)
	}

	resp, err := c.httpClient.Post(
		c.baseURL+endpoint,
		"application/json",
		bytes.NewReader(body),
	)
	if err != nil {
		return fmt.Errorf("http post: %w", err)
	}
	defer resp.Body.Close()

	var rpcResp rpcResponse
	if err := json.NewDecoder(resp.Body).Decode(&rpcResp); err != nil {
		return fmt.Errorf("decode rpc response: %w", err)
	}

	if rpcResp.Error != nil {
		return &RPCError{
			Code:    rpcResp.Error.Code,
			Message: rpcResp.Error.Message,
		}
	}

	if out != nil {
		if err := json.Unmarshal(rpcResp.Result, out); err != nil {
			return fmt.Errorf("decode rpc result: %w", err)
		}
	}

	return nil
}
