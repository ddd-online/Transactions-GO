package mcp

import (
	"encoding/json"
)

// JSON-RPC 2.0 request
type JsonRpcRequest struct {
	JsonRpc string          `json:"jsonrpc"`
	ID      interface{}     `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

// JSON-RPC 2.0 response
type JsonRpcResponse struct {
	JsonRpc string        `json:"jsonrpc"`
	ID      interface{}   `json:"id"`
	Result  interface{}   `json:"result,omitempty"`
	Error   *JsonRpcError `json:"error,omitempty"`
}

type JsonRpcError struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data,omitempty"`
}

func (e *JsonRpcError) Error() string {
	return e.Message
}

// JSON-RPC 2.0 error codes
const (
	CodeParseError     = -32700
	CodeInvalidRequest = -32600
	CodeMethodNotFound = -32601
	CodeInvalidParams  = -32602
	CodeInternalError  = -32603
)

// ParseRequest 解析 JSON-RPC 请求
func ParseRequest(data []byte) (*JsonRpcRequest, error) {
	var req JsonRpcRequest
	if err := json.Unmarshal(data, &req); err != nil {
		return nil, err
	}
	if req.JsonRpc != "2.0" {
		return nil, &JsonRpcError{Code: CodeInvalidRequest, Message: "Invalid JSON-RPC version"}
	}
	return &req, nil
}

// NewSuccessResponse 创建成功响应
func NewSuccessResponse(id interface{}, result interface{}) JsonRpcResponse {
	return JsonRpcResponse{
		JsonRpc: "2.0",
		ID:      id,
		Result:  result,
	}
}

// NewErrorResponse 创建错误响应
func NewErrorResponse(id interface{}, code int, message string) JsonRpcResponse {
	return JsonRpcResponse{
		JsonRpc: "2.0",
		ID:      id,
		Error: &JsonRpcError{
			Code:    code,
			Message: message,
		},
	}
}
