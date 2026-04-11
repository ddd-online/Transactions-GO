package mcp

import (
	"sync"
)

var (
	mcpServer   *McpServer
	mcpServerMu sync.Mutex
)

func GetMcpServer() *McpServer {
	mcpServerMu.Lock()
	defer mcpServerMu.Unlock()
	if mcpServer == nil {
		mcpServer = NewMcpServer()
	}
	return mcpServer
}

func StartMcpServer() error {
	server := GetMcpServer()
	return server.Start()
}

func StopMcpServer() error {
	server := GetMcpServer()
	return server.Stop()
}

func GetMcpStatus() map[string]interface{} {
	server := GetMcpServer()
	return map[string]interface{}{
		"running": server.IsRunning(),
		"port":    ServerPort,
	}
}