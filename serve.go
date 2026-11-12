// Reference stdio server for extensions. It implements the same NDJSON
// JSON-RPC protocol used by the host and can be copied into a Go plugin SDK.
package extensionapi

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// Serve runs an Extension over stdin/stdout until stdin closes or a shutdown
// request arrives. stdout is reserved for protocol messages; logs go to stderr.
func Serve(ext Extension) {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 64*1024), 16*1024*1024)
	out := json.NewEncoder(os.Stdout)

	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		var req Request
		if err := json.Unmarshal(line, &req); err != nil {
			if err := out.Encode(Response{
				JSONRPC: JSONRPCVersion,
				Error:   &RPCError{Code: -32700, Message: "invalid JSON: " + err.Error()},
			}); err != nil {
				log.Printf("ERROR [extension serve]: %v", err)
				return
			}
			continue
		}

		result, rpcErr := dispatch(ext, req)
		resp := Response{JSONRPC: JSONRPCVersion, ID: req.ID, Error: rpcErr}
		if rpcErr == nil {
			data, err := json.Marshal(result)
			if err != nil {
				resp.Error = &RPCError{Code: -32603, Message: err.Error()}
			} else {
				resp.Result = data
			}
		}
		if err := out.Encode(resp); err != nil {
			log.Printf("ERROR [extension serve]: %v", err)
			return
		}
		if req.Method == "shutdown" {
			return
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("ERROR [extension serve]: %v", err)
	}
}

func dispatch(ext Extension, req Request) (any, *RPCError) {
	ctx := context.Background()
	switch req.Method {
	case "info":
		info, err := ext.Info(ctx)
		if err != nil {
			return nil, &RPCError{Code: -32000, Message: err.Error()}
		}
		return info, nil
	case "contribute":
		var in ContributionRequest
		if err := json.Unmarshal(req.Params, &in); err != nil {
			return nil, &RPCError{Code: -32602, Message: fmt.Sprintf("invalid contribute params: %v", err)}
		}
		out, err := ext.Contribute(ctx, &in)
		if err != nil {
			return nil, &RPCError{Code: -32000, Message: err.Error()}
		}
		return out, nil
	case "handle":
		var in Event
		if err := json.Unmarshal(req.Params, &in); err != nil {
			return nil, &RPCError{Code: -32602, Message: fmt.Sprintf("invalid handle params: %v", err)}
		}
		out, err := ext.Handle(ctx, &in)
		if err != nil {
			return nil, &RPCError{Code: -32000, Message: err.Error()}
		}
		return out, nil
	case "shutdown":
		return map[string]bool{"ok": true}, nil
	default:
		return nil, &RPCError{Code: -32601, Message: "unknown method: " + req.Method}
	}
}
