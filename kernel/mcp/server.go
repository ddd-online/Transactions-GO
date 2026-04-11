package mcp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"sync"

	"github.com/sirupsen/logrus"
)

// Protocol version and server configuration
const (
	ProtocolVersion = "2024-11-05"
	ServerPort      = 31944
)

// McpServer is the main MCP TCP server
type McpServer struct {
	listener net.Listener
	wg       sync.WaitGroup
	tools    *ToolHandler
	mu       sync.Mutex
	running  bool
}

// NewMcpServer creates a new MCP server instance
func NewMcpServer() *McpServer {
	return &McpServer{
		tools: NewToolHandler(),
	}
}

// Start starts the TCP listener and begins accepting connections
func (s *McpServer) Start() error {
	s.mu.Lock()
	if s.running {
		s.mu.Unlock()
		return fmt.Errorf("server is already running")
	}

	addr := fmt.Sprintf("127.0.0.1:%d", ServerPort)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		s.mu.Unlock()
		return fmt.Errorf("failed to listen on %s: %w", addr, err)
	}

	s.listener = listener
	s.running = true
	s.mu.Unlock()

	logrus.Infof("MCP server listening on %s", addr)

	s.wg.Add(1)
	go s.acceptConnections()

	return nil
}

// Stop stops the server and closes the listener
func (s *McpServer) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.running {
		return fmt.Errorf("server is not running")
	}

	s.running = false

	if s.listener != nil {
		if err := s.listener.Close(); err != nil {
			return fmt.Errorf("failed to close listener: %w", err)
		}
	}

	s.wg.Wait()
	logrus.Info("MCP server stopped")
	return nil
}

// IsRunning returns the running status of the server
func (s *McpServer) IsRunning() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.running
}

// acceptConnections accepts incoming TCP connections
func (s *McpServer) acceptConnections() {
	defer s.wg.Done()

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			s.mu.Lock()
			running := s.running
			s.mu.Unlock()

			if !running {
				// Server was stopped, exit gracefully
				return
			}
			logrus.Errorf("failed to accept connection: %v", err)
			continue
		}

		s.wg.Add(1)
		go s.handleConnection(conn)
	}
}

// handleConnection handles a single TCP connection
func (s *McpServer) handleConnection(conn net.Conn) {
	defer s.wg.Done()
	defer conn.Close()

	logrus.Debugf("new connection from %s", conn.RemoteAddr())

	reader := bufio.NewReader(conn)
	encoder := json.NewEncoder(conn)

	// Send initialized notification (no id field per JSON-RPC 2.0 spec)
	initialized := JsonRpcNotification{
		JsonRpc: "2.0",
		Method:  "notifications/initialized",
		Params: map[string]interface{}{
			"protocolVersion": ProtocolVersion,
		},
	}
	if err := encoder.Encode(initialized); err != nil {
		logrus.Errorf("failed to send initialized notification: %v", err)
		return
	}

	// Handle requests in a loop
	for {
		data, err := reader.ReadBytes('\n')
		if err != nil {
			if err.Error() != "EOF" {
				logrus.Errorf("read error: %v", err)
			}
			return
		}

		// Remove trailing newline
		if len(data) > 0 && data[len(data)-1] == '\n' {
			data = data[:len(data)-1]
		}
		// Also handle carriage return
		if len(data) > 0 && data[len(data)-1] == '\r' {
			data = data[:len(data)-1]
		}

		if len(data) == 0 {
			continue
		}

		req, err := ParseRequest(data)
		if err != nil {
			logrus.Errorf("failed to parse request: %v", err)
			resp := NewErrorResponse(nil, CodeInvalidRequest, err.Error())
			encoder.Encode(resp)
			continue
		}

		resp := s.handleRequest(conn, encoder, req)
		if resp != nil {
			if err := encoder.Encode(resp); err != nil {
				logrus.Errorf("failed to encode response: %v", err)
			}
		}
	}
}

// handleRequest routes requests to the appropriate handler
func (s *McpServer) handleRequest(conn net.Conn, encoder *json.Encoder, req *JsonRpcRequest) *JsonRpcResponse {
	switch req.Method {
	case "initialize":
		resp := s.handleInitialize(req)
		return &resp
	case "tools/list":
		resp := s.handleToolsList(req)
		return &resp
	case "tools/call":
		resp := s.handleToolsCall(req)
		return &resp
	default:
		logrus.Warnf("unknown method: %s", req.Method)
		resp := NewErrorResponse(req.ID, CodeMethodNotFound, fmt.Sprintf("method not found: %s", req.Method))
		return &resp
	}
}

// handleInitialize handles the initialize request
func (s *McpServer) handleInitialize(req *JsonRpcRequest) JsonRpcResponse {
	return NewSuccessResponse(req.ID, map[string]interface{}{
		"protocolVersion": ProtocolVersion,
		"capabilities": map[string]interface{}{
			"tools": struct{}{},
		},
		"serverInfo": map[string]interface{}{
			"name":    "billadm-mcp",
			"version": "1.0.0",
		},
	})
}

// handleToolsList handles the tools/list request
func (s *McpServer) handleToolsList(req *JsonRpcRequest) JsonRpcResponse {
	tools := s.tools.ListTools()
	return NewSuccessResponse(req.ID, map[string]interface{}{
		"tools": tools,
	})
}

// handleToolsCall handles the tools/call request
func (s *McpServer) handleToolsCall(req *JsonRpcRequest) JsonRpcResponse {
	var params struct {
		Name      string          `json:"name"`
		Arguments json.RawMessage `json:"arguments,omitempty"`
	}

	if err := json.Unmarshal(req.Params, &params); err != nil {
		return NewErrorResponse(req.ID, CodeInvalidParams, fmt.Sprintf("invalid params: %v", err))
	}

	result, err := s.tools.CallTool(params.Name, params.Arguments)
	if err != nil {
		return NewErrorResponse(req.ID, CodeInternalError, err.Error())
	}

	return NewSuccessResponse(req.ID, map[string]interface{}{
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": fmt.Sprintf("%v", result),
			},
		},
	})
}

