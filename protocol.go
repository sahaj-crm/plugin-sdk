// Minimal JSON-RPC wire types shared by the host and the reference stdio
// server. External SDKs implement the same one-object-per-line protocol.
package extensionapi

import "encoding/json"

const JSONRPCVersion = "2.0"

// Request is one line received by a plugin.
type Request struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      int             `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

// Response is one line written by a plugin.
type Response struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      int             `json:"id"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *RPCError       `json:"error,omitempty"`
}

// RPCError is a protocol-level error.
type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
