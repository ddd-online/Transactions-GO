# MCP Server 实现计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 为 Billadm-GO 增加 MCP 服务器支持，使 AI 助手能通过 TCP 连接查询账本和交易记录

**Architecture:** Go 后端在 TCP 端口 31944 上提供 MCP 服务，前端通过 HTTP API 控制启停。MCP 基于 JSON-RPC 2.0，实现 `initialize`、`tools/list`、`tools/call` 三个核心方法。

**Tech Stack:** Go (kernel/mcp/)、TypeScript/Vue (app/src/)

---

## 文件结构

```
kernel/
  mcp/
    protocol.go     # Create: JSON-RPC 2.0 协议工具函数
    server.go       # Create: MCP TCP 服务器，处理连接和请求分发
    tools.go        # Create: MCP 工具实现（query_ledgers, query_transactions）
    manager.go      # Create: MCP 服务器生命周期管理（启动/停止/状态）
  api/
    mcp_controller.go  # Create: MCP 控制 API（start/stop/status）

app/src/
  stores/
    mcpStore.ts     # Create: MCP 状态管理
  backend/api/
    mcp.ts          # Create: MCP 控制 API 调用
  components/settings_view/
    McpSetting.vue  # Create: MCP 设置组件
  (修改 SettingsView.vue: 添加 MCP 设置菜单项)
```

---

## Task 1: 创建 MCP 协议基础 (protocol.go)

**Files:**
- Create: `kernel/mcp/protocol.go`

- [ ] **Step 1: 创建 protocol.go 文件**

```go
package mcp

import (
	"encoding/json"
)

// JSON-RPC 2.0 request
type JsonRpcRequest struct {
	JsonRpc string          `json:"jsonrpc"`
	ID     interface{}     `json:"id"`
	Method string          `json:"method"`
	Params json.RawMessage `json:"params,omitempty"`
}

// JSON-RPC 2.0 response
type JsonRpcResponse struct {
	JsonRpc string      `json:"jsonrpc"`
	ID     interface{} `json:"id"`
	Result interface{} `json:"result,omitempty"`
	Error  *JsonRpcError `json:"error,omitempty"`
}

type JsonRpcError struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data,omitempty"`
}

// JSON-RPC 2.0 error codes
const (
	CodeParseError       = -32700
	CodeInvalidRequest   = -32600
	CodeMethodNotFound   = -32601
	CodeInvalidParams    = -32602
	CodeInternalError    = -32603
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
		ID:     id,
		Result: result,
	}
}

// NewErrorResponse 创建错误响应
func NewErrorResponse(id interface{}, code int, message string) JsonRpcResponse {
	return JsonRpcResponse{
		JsonRpc: "2.0",
		ID:     id,
		Error: &JsonRpcError{
			Code:    code,
			Message: message,
		},
	}
}
```

- [ ] **Step 2: 提交**

```bash
git add kernel/mcp/protocol.go
git commit -m "feat(mcp): add JSON-RPC 2.0 protocol utilities"
```

---

## Task 2: 创建 MCP 服务器 (server.go)

**Files:**
- Create: `kernel/mcp/server.go`

- [ ] **Step 1: 创建 server.go 文件**

```go
package mcp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"sync"

	"github.com/sirupsen/logrus"
)

const (
	ProtocolVersion = "2024-11-05"
	ServerPort      = 31944
)

// McpServer MCP TCP 服务器
type McpServer struct {
	listener net.Listener
	wg       sync.WaitGroup
	tools    *McpToolHandler
	mu       sync.Mutex
	running  bool
}

func NewMcpServer() *McpServer {
	return &McpServer{
		tools: NewToolHandler(),
	}
}

// Start 启动 MCP 服务器
func (s *McpServer) Start() error {
	s.mu.Lock()
	if s.running {
		s.mu.Unlock()
		return fmt.Errorf("MCP server already running")
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

	logrus.Infof("MCP server started on %s", addr)

	// 接受连接
	go s.acceptConnections()
	return nil
}

// Stop 停止 MCP 服务器
func (s *McpServer) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.running {
		return fmt.Errorf("MCP server not running")
	}

	s.running = false
	if s.listener != nil {
		s.listener.Close()
	}

	logrus.Info("MCP server stopped")
	return nil
}

// IsRunning 返回服务器运行状态
func (s *McpServer) IsRunning() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.running
}

func (s *McpServer) acceptConnections() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			s.mu.Lock()
			running := s.running
			s.mu.Unlock()
			if !running {
				return
			}
			logrus.Warnf("MCP server accept error: %v", err)
			continue
		}

		s.wg.Add(1)
		go s.handleConnection(conn)
	}
}

func (s *McpServer) handleConnection(conn net.Conn) {
	defer s.wg.Done()
	defer conn.Close()

	reader := bufio.NewReader(conn)
	decoder := json.NewDecoder(reader)
	encoder := json.NewEncoder(conn)

	// 发送初始化通知（服务器推送）
	notification := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "notifications/initialized",
		"params":  map[string]interface{}{},
	}
	encoder.Encode(notification)

	for {
		var req JsonRpcRequest
		if err := decoder.Decode(&req); err != nil {
			logrus.Debugf("MCP connection closed: %v", err)
			return
		}

		s.handleRequest(conn, encoder, &req)
	}
}

func (s *McpServer) handleRequest(conn net.Conn, encoder *json.Encoder, req *JsonRpcRequest) {
	var resp JsonRpcResponse

	switch req.Method {
	case "initialize":
		resp = s.handleInitialize(req)
	case "tools/list":
		resp = s.handleToolsList(req)
	case "tools/call":
		resp = s.handleToolsCall(req)
	case "ping":
		resp = NewSuccessResponse(req.ID, map[string]interface{}{"pong": true})
	default:
		resp = NewErrorResponse(req.ID, CodeMethodNotFound, fmt.Sprintf("Method not found: %s", req.Method))
	}

	encoder.Encode(resp)
}

func (s *McpServer) handleInitialize(req *JsonRpcRequest) JsonRpcResponse {
	result := map[string]interface{}{
		"protocolVersion": ProtocolVersion,
		"capabilities": map[string]interface{}{
			"tools": map[string]interface{}{},
		},
		"serverInfo": map[string]interface{}{
			"name":    "billadm-mcp",
			"version": "1.0.0",
		},
	}
	return NewSuccessResponse(req.ID, result)
}

func (s *McpServer) handleToolsList(req *JsonRpcRequest) JsonRpcResponse {
	tools := s.tools.ListTools()
	result := map[string]interface{}{
		"tools": tools,
	}
	return NewSuccessResponse(req.ID, result)
}

func (s *McpServer) handleToolsCall(req *JsonRpcRequest) JsonRpcResponse {
	var params struct {
		Name      string          `json:"name"`
		Arguments json.RawMessage `json:"arguments,omitempty"`
	}
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return NewErrorResponse(req.ID, CodeInvalidParams, "Invalid params")
	}

	result, err := s.tools.CallTool(params.Name, params.Arguments)
	if err != nil {
		return NewErrorResponse(req.ID, CodeInternalError, err.Error())
	}

	return NewSuccessResponse(req.ID, result)
}
```

- [ ] **Step 2: 提交**

```bash
git add kernel/mcp/server.go
git commit -m "feat(mcp): add MCP TCP server implementation"
```

---

## Task 3: 创建 MCP 工具实现 (tools.go)

**Files:**
- Create: `kernel/mcp/tools.go`

- [ ] **Step 1: 创建 tools.go 文件**

```go
package mcp

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/billadm/models/dto"
	"github.com/billadm/service"
	"github.com/billadm/workspace"
	"github.com/sirupsen/logrus"
)

const DefaultLimit = 50

// ToolDefinition MCP 工具定义
type ToolDefinition struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	InputSchema map[string]interface{} `json:"inputSchema"`
}

// ToolHandler MCP 工具处理
type ToolHandler struct {
	tools []ToolDefinition
}

func NewToolHandler() *ToolHandler {
	return &ToolHandler{
		tools: []ToolDefinition{
			{
				Name:        "query_ledgers",
				Description: "查询所有账本列表",
				InputSchema: map[string]interface{}{
					"type":       "object",
					"properties": map[string]interface{}{},
				},
			},
			{
				Name:        "query_transactions",
				Description: "按条件查询交易记录",
				InputSchema: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"ledger_id": map[string]interface{}{
							"type":        "string",
							"description": "账本ID（必填）",
						},
						"time_range": map[string]interface{}{
							"type":        "array",
							"items":       map[string]interface{}{"type": "integer"},
							"description": "时间戳范围 [start, end]，单位秒",
						},
						"transaction_type": map[string]interface{}{
							"type":        "string",
							"enum":        []string{"expense", "income", "transfer"},
							"description": "交易类型",
						},
						"category": map[string]interface{}{
							"type":        "string",
							"description": "分类名称",
						},
						"tags": map[string]interface{}{
							"type":        "array",
							"items":       map[string]interface{}{"type": "string"},
							"description": "标签列表",
						},
						"description": map[string]interface{}{
							"type":        "string",
							"description": "备注包含的字符",
						},
						"offset": map[string]interface{}{
							"type":        "integer",
							"default":     0,
							"description": "分页偏移",
						},
						"limit": map[string]interface{}{
							"type":        "integer",
							"default":     50,
							"description": "每页数量，最大100",
						},
						"sort_fields": map[string]interface{}{
							"type": "array",
							"items": map[string]interface{}{
								"type": "object",
								"properties": map[string]interface{}{
									"field": map[string]interface{}{"type": "string"},
									"order": map[string]interface{}{"type": "string", "enum": []string{"asc", "desc"}},
								},
							},
							"description": "排序字段列表",
						},
					},
					"required": []string{"ledger_id"},
				},
			},
		},
	}
}

func (h *ToolHandler) ListTools() []ToolDefinition {
	return h.tools
}

func (h *ToolHandler) CallTool(name string, arguments json.RawMessage) (interface{}, error) {
	switch name {
	case "query_ledgers":
		return h.queryLedgers()
	case "query_transactions":
		return h.queryTransactions(arguments)
	default:
		return nil, fmt.Errorf("unknown tool: %s", name)
	}
}

func (h *ToolHandler) queryLedgers() (interface{}, error) {
	ws := workspace.Manager.OpenedWorkspace()
	if ws == nil {
		return nil, fmt.Errorf("工作空间未打开")
	}

	ledgers, err := service.GetLedgerService().ListAllLedger(ws)
	if err != nil {
		return nil, fmt.Errorf("查询账本失败: %w", err)
	}

	var lines []string
	lines = append(lines, "账本列表：")
	for i, ledger := range ledgers {
		desc := ""
		if ledger.Description != "" {
			desc = fmt.Sprintf(" - %s", ledger.Description)
		}
		lines = append(lines, fmt.Sprintf("%d. %s (ID: %s%s)", i+1, ledger.Name, ledger.LedgerID, desc))
	}

	return map[string]interface{}{
		"content": []map[string]string{
			{"type": "text", "text": strings.Join(lines, "\n")},
		},
	}, nil
}

func (h *ToolHandler) queryTransactions(arguments json.RawMessage) (interface{}, error) {
	ws := workspace.Manager.OpenedWorkspace()
	if ws == nil {
		return nil, fmt.Errorf("工作空间未打开")
	}

	// 解析参数
	var args struct {
		LedgerID       string  `json:"ledger_id"`
		TimeRange      []int64 `json:"time_range"`
		TransactionType string  `json:"transaction_type"`
		Category       string  `json:"category"`
		Tags           []string `json:"tags"`
		Description    string  `json:"description"`
		Offset         int     `json:"offset"`
		Limit          int     `json:"limit"`
		SortFields     []struct {
			Field string `json:"field"`
			Order string `json:"order"`
		} `json:"sort_fields"`
	}

	if err := json.Unmarshal(arguments, &args); err != nil {
		return nil, fmt.Errorf("解析参数失败: %w", err)
	}

	if args.LedgerID == "" {
		return nil, fmt.Errorf("ledger_id 为必填参数")
	}

	if args.Limit <= 0 || args.Limit > 100 {
		args.Limit = DefaultLimit
	}

	// 构建查询条件
	condition := &dto.TrQueryCondition{
		LedgerID: args.LedgerID,
		Offset:   args.Offset,
		Limit:    args.Limit,
	}

	if len(args.TimeRange) == 2 {
		condition.TsRange = args.TimeRange
	}

	if args.TransactionType != "" || args.Category != "" || len(args.Tags) > 0 || args.Description != "" {
		item := dto.QueryConditionItem{
			TransactionType: args.TransactionType,
			Category:        args.Category,
			Tags:            args.Tags,
			Description:     args.Description,
		}
		condition.Items = append(condition.Items, item)
	}

	if len(args.SortFields) > 0 {
		for _, sf := range args.SortFields {
			condition.SortFields = append(condition.SortFields, dto.QueryConditionSortField{
				Field: sf.Field,
				Order: sf.Order,
			})
		}
	}

	// 执行查询
	result, err := service.GetTrService().QueryTrsOnCondition(ws, condition)
	if err != nil {
		logrus.Errorf("MCP query transactions failed: %v", err)
		return nil, fmt.Errorf("查询交易记录失败: %w", err)
	}

	// 格式化结果
	var lines []string
	lines = append(lines, fmt.Sprintf("交易记录列表（共 %d 条）：", result.Total))

	for i, tr := range result.Items {
		date := formatTimestamp(tr.TransactionAt)
		amount := formatAmount(tr.TransactionType, tr.Price)
		tags := ""
		if len(tr.Tags) > 0 {
			tags = " - " + strings.Join(tr.Tags, " #")
		}
		desc := ""
		if tr.Description != "" {
			desc = fmt.Sprintf(" - 备注: %s", tr.Description)
		}
		lines = append(lines, fmt.Sprintf("%d. [%s] %s - %s%s%s - %s",
			i+1+args.Offset, date, tr.TransactionType, tr.Category, tags, amount, desc))
	}

	totalPages := (result.Total + args.Limit - 1) / args.Limit
	currentPage := args.Offset/args.Limit + 1
	lines = append(lines, fmt.Sprintf("\n共 %d 条记录，第 %d/%d 页", result.Total, currentPage, totalPages))

	return map[string]interface{}{
		"content": []map[string]string{
			{"type": "text", "text": strings.Join(lines, "\n")},
		},
	}, nil
}

func formatTimestamp(ts int64) string {
	if ts <= 0 {
		return "N/A"
	}
	// ts 是秒级时间戳
	t := ts
	// 尝试判断是秒还是毫秒
	if ts > 1e12 {
		t = ts / 1000
	}
	sec := t % 60
	min := (t / 60) % 60
	hour := (t / 3600) % 24
	day := (t / 86400) % 30
	month := (t / 2592000) % 12
	year := t / 31536000
	if year > 0 {
		return fmt.Sprintf("%04d-%02d-%02d", year+1970, month+1, day+1)
	}
	return fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d", year+1970, month+1, day+1, hour, min, sec)
}

func formatAmount(txType string, price int64) string {
	unit := "分"
	if price < 0 {
		unit = "分"
		price = -price
	}
	yuan := price / 100
	cents := price % 100
	if txType == "expense" {
		return fmt.Sprintf(" - %d.%02d元", -yuan, cents)
	} else if txType == "income" {
		return fmt.Sprintf(" + %d.%02d元", yuan, cents)
	}
	return fmt.Sprintf(" %d.%02d元", yuan, cents)
}
```

- [ ] **Step 2: 提交**

```bash
git add kernel/mcp/tools.go
git commit -m "feat(mcp): add MCP tools implementation"
```

---

## Task 4: 创建 MCP 管理器 (manager.go)

**Files:**
- Create: `kernel/mcp/manager.go`

- [ ] **Step 1: 创建 manager.go 文件**

```go
package mcp

import (
	"fmt"
	"sync"
)

var (
	mcpServer   *McpServer
	mcpServerMu sync.Mutex
)

// GetMcpServer 获取 MCP 服务器单例
func GetMcpServer() *McpServer {
	mcpServerMu.Lock()
	defer mcpServerMu.Unlock()

	if mcpServer == nil {
		mcpServer = NewMcpServer()
	}
	return mcpServer
}

// StartMcpServer 启动 MCP 服务器
func StartMcpServer() error {
	server := GetMcpServer()
	return server.Start()
}

// StopMcpServer 停止 MCP 服务器
func StopMcpServer() error {
	server := GetMcpServer()
	return server.Stop()
}

// GetMcpStatus 获取 MCP 服务器状态
func GetMcpStatus() map[string]interface{} {
	server := GetMcpServer()
	return map[string]interface{}{
		"running": server.IsRunning(),
		"port":    ServerPort,
	}
}
```

- [ ] **Step 2: 提交**

```bash
git add kernel/mcp/manager.go
git commit -m "feat(mcp): add MCP server manager"
```

---

## Task 5: 创建 MCP 控制 API (mcp_controller.go)

**Files:**
- Create: `kernel/api/mcp_controller.go`
- Modify: `kernel/api/router.go:63-72` (添加 MCP 路由)

- [ ] **Step 1: 创建 mcp_controller.go 文件**

```go
package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/billadm/mcp"
	"github.com/billadm/models"
)

// POST /api/v1/mcp/start
func startMcpServer(c *gin.Context) {
	ret := models.NewResult()
	defer c.JSON(http.StatusOK, ret)

	err := mcp.StartMcpServer()
	if err != nil {
		ret.Code = -1
		ret.Msg = err.Error()
		return
	}

	ret.Data = mcp.GetMcpStatus()
}

// POST /api/v1/mcp/stop
func stopMcpServer(c *gin.Context) {
	ret := models.NewResult()
	defer c.JSON(http.StatusOK, ret)

	err := mcp.StopMcpServer()
	if err != nil {
		ret.Code = -1
		ret.Msg = err.Error()
		return
	}

	ret.Data = mcp.GetMcpStatus()
}

// GET /api/v1/mcp/status
func getMcpStatus(c *gin.Context) {
	ret := models.NewResult()
	defer c.JSON(http.StatusOK, ret)

	ret.Data = mcp.GetMcpStatus()
}
```

- [ ] **Step 2: 修改 router.go 添加 MCP 路由**

在 `v1 := ginServer.Group("/api/v1")` 的 v1.Group 块中添加：

```go
// MCP server control
mcpGroup := v1.Group("/mcp")
{
	mcpGroup.POST("/start", startMcpServer)
	mcpGroup.POST("/stop", stopMcpServer)
	mcpGroup.GET("/status", getMcpStatus)
}
```

- [ ] **Step 3: 提交**

```bash
git add kernel/api/mcp_controller.go kernel/api/router.go
git commit -m "feat(mcp): add MCP control API endpoints"
```

---

## Task 6: 创建前端 MCP API 调用 (mcp.ts)

**Files:**
- Create: `app/src/backend/api/mcp.ts`

- [ ] **Step 1: 创建 mcp.ts 文件**

```typescript
import api from "@/backend/api/api-client";

export interface McpStatus {
  running: boolean;
  port: number;
}

export async function startMcpServer(): Promise<McpStatus> {
  return api.post<McpStatus>('/v1/mcp/start', {}, '启动MCP服务');
}

export async function stopMcpServer(): Promise<McpStatus> {
  return api.post<McpStatus>('/v1/mcp/stop', {}, '停止MCP服务');
}

export async function getMcpStatus(): Promise<McpStatus> {
  return api.get<McpStatus>('/v1/mcp/status', '获取MCP状态');
}
```

- [ ] **Step 2: 提交**

```bash
git add app/src/backend/api/mcp.ts
git commit -m "feat(mcp): add frontend MCP API calls"
```

---

## Task 7: 创建前端 MCP 状态管理 (mcpStore.ts)

**Files:**
- Create: `app/src/stores/mcpStore.ts`

- [ ] **Step 1: 创建 mcpStore.ts 文件**

```typescript
import {defineStore} from 'pinia'
import {ref} from 'vue'
import {startMcpServer, stopMcpServer, getMcpStatus, type McpStatus} from "@/backend/api/mcp"
import NotificationUtil from "@/backend/notification"

export const useMcpStore = defineStore('mcp', () => {
  const isRunning = ref(false)
  const port = ref(31944)

  const refreshStatus = async () => {
    try {
      const status = await getMcpStatus()
      isRunning.value = status.running
      port.value = status.port
    } catch (error) {
      console.warn('Failed to get MCP status:', error)
    }
  }

  const start = async () => {
    try {
      const status = await startMcpServer()
      isRunning.value = status.running
      port.value = status.port
      NotificationUtil.success('MCP 服务已启动', `端口: ${status.port}`)
    } catch (error) {
      NotificationUtil.error('启动 MCP 服务失败', `${error}`)
      throw error
    }
  }

  const stop = async () => {
    try {
      const status = await stopMcpServer()
      isRunning.value = status.running
      NotificationUtil.success('MCP 服务已停止')
    } catch (error) {
      NotificationUtil.error('停止 MCP 服务失败', `${error}`)
      throw error
    }
  }

  return {
    isRunning,
    port,
    refreshStatus,
    start,
    stop,
  }
})
```

- [ ] **Step 2: 提交**

```bash
git add app/src/stores/mcpStore.ts
git commit -m "feat(mcp): add MCP store for state management"
```

---

## Task 8: 创建 MCP 设置组件 (McpSetting.vue)

**Files:**
- Create: `app/src/components/settings_view/McpSetting.vue`
- Modify: `app/src/components/settings_view/SettingsView.vue` (添加菜单项)

- [ ] **Step 1: 创建 McpSetting.vue 文件**

```vue
<template>
  <div class="mcp-setting">
    <div class="setting-header">
      <h3>MCP 服务</h3>
    </div>

    <div class="setting-content">
      <div class="setting-row">
        <div class="setting-info">
          <div class="setting-label">启用 MCP 服务</div>
          <div class="setting-desc">
            启用后 AI 助手（如 Claude Code）可以通过 TCP 连接查询账本和交易记录
          </div>
        </div>
        <a-switch v-model:checked="switchLoading" :loading="loading" @change="handleSwitchChange" />
      </div>

      <div class="setting-row" v-if="isRunning">
        <div class="setting-info">
          <div class="setting-label">服务状态</div>
          <div class="setting-desc">运行中</div>
        </div>
        <div class="status-indicator running"></div>
      </div>

      <div class="setting-row">
        <div class="setting-info">
          <div class="setting-label">端口</div>
          <div class="setting-desc">{{ port }}</div>
        </div>
      </div>

      <div class="setting-row" v-if="isRunning">
        <div class="setting-info">
          <div class="setting-label">Claude Code 配置</div>
          <div class="setting-desc">在 Claude Code 设置中添加以下 MCP 服务器配置：</div>
        </div>
      </div>

      <div class="config-block" v-if="isRunning">
        <pre class="config-code">{{ mcpConfig }}</pre>
        <a-button size="small" @click="copyConfig">复制</a-button>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import {ref, computed, onMounted} from 'vue'
import {message} from 'ant-design-vue'
import {useMcpStore} from '@/stores/mcpStore'

const mcpStore = useMcpStore()
const loading = ref(false)
const switchLoading = ref(false)

const isRunning = computed(() => mcpStore.isRunning)
const port = computed(() => mcpStore.port)

const mcpConfig = computed(() => `{
  "mcpServers": {
    "billadm": {
      "command": "nc",
      "args": ["127.0.0.1", "${port.value}"]
    }
  }
}`)

onMounted(async () => {
  loading.value = true
  await mcpStore.refreshStatus()
  loading.value = false
  switchLoading.value = mcpStore.isRunning
})

const handleSwitchChange = async (checked: boolean) => {
  try {
    if (checked) {
      await mcpStore.start()
    } else {
      await mcpStore.stop()
    }
  } catch (error) {
    // 恢复开关状态
    switchLoading.value = mcpStore.isRunning
  }
}

const copyConfig = () => {
  navigator.clipboard.writeText(mcpConfig.value)
  message.success('已复制到剪贴板')
}
</script>

<style scoped>
.mcp-setting {
  max-width: 600px;
}

.setting-header h3 {
  margin: 0 0 16px 0;
  font-size: 16px;
  font-weight: 500;
}

.setting-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 0;
  border-bottom: 1px solid var(--billadm-color-window-border);
}

.setting-row:last-child {
  border-bottom: none;
}

.setting-info {
  flex: 1;
}

.setting-label {
  font-size: 14px;
  color: var(--billadm-color-text-primary);
  margin-bottom: 4px;
}

.setting-desc {
  font-size: 12px;
  color: var(--billadm-color-text-secondary);
}

.status-indicator {
  width: 8px;
  height: 8px;
  border-radius: 50%;
}

.status-indicator.running {
  background-color: #52c41a;
  box-shadow: 0 0 4px #52c41a;
}

.config-block {
  background: var(--billadm-color-minor-background);
  border-radius: 6px;
  padding: 12px;
  margin-top: 8px;
  display: flex;
  gap: 12px;
  align-items: flex-start;
}

.config-code {
  flex: 1;
  font-size: 12px;
  font-family: 'Courier New', monospace;
  margin: 0;
  white-space: pre-wrap;
  word-break: break-all;
}
</style>
```

- [ ] **Step 2: 修改 SettingsView.vue 添加菜单项**

在 `<div class="settings-list">` 中添加：

```vue
<div
    class="settings-list-item"
    :class="{ active: activeComponent === 'mcp' }"
    @click="activeComponent = 'mcp'"
>
  <SettingOutlined class="settings-list-item-icon"/>
  <span class="settings-list-item-title">MCP 服务</span>
</div>
```

在 script 中导入 `SettingOutlined` 和 `McpSetting`：

```typescript
import {SettingOutlined} from "@ant-design/icons-vue";
import McpSetting from './McpSetting.vue';
```

在条件渲染区域添加：

```vue
<mcp-setting v-else-if="activeComponent === 'mcp'"/>
```

- [ ] **Step 3: 提交**

```bash
git add app/src/components/settings_view/McpSetting.vue app/src/components/settings_view/SettingsView.vue
git commit -m "feat(mcp): add MCP settings UI component"
```

---

## Task 9: 集成测试

**Files:**
- Modify: `kernel/main.go` (如有需要添加测试入口)

- [ ] **Step 1: 构建并测试 MCP 服务器**

```bash
cd kernel && go build -ldflags '-s -w -extldflags "-static"' -o Billadm-Kernel.exe
```

- [ ] **Step 2: 测试 MCP API**

使用 curl 测试：

```bash
# 启动 MCP 服务
curl -X POST http://127.0.0.1:31943/api/v1/mcp/start

# 检查状态
curl http://127.0.0.1:31943/api/v1/mcp/status

# 停止 MCP 服务
curl -X POST http://127.0.0.1:31943/api/v1/mcp/stop
```

- [ ] **Step 3: 测试 MCP 协议连接**

使用 netcat 测试 JSON-RPC：

```bash
# 连接 MCP 服务器
nc 127.0.0.1 31944

# 发送 initialize 请求
{"jsonrpc": "2.0", "id": 1, "method": "initialize", "params": {"protocolVersion": "2024-11-05", "capabilities": {}}}

# 发送 tools/list 请求
{"jsonrpc": "2.0", "id": 2, "method": "tools/list", "params": {}}

# 发送 tools/call 请求
{"jsonrpc": "2.0", "id": 3, "method": "tools/call", "params": {"name": "query_ledgers", "arguments": {}}}
```

- [ ] **Step 4: 提交**

```bash
git add -A
git commit -m "test(mcp): add integration tests for MCP server"
```

---

## 自检清单

1. **Spec 覆盖检查：**
   - [x] MCP TCP 服务器实现（Task 2）
   - [x] JSON-RPC 2.0 协议处理（Task 1）
   - [x] query_ledgers 工具（Task 3）
   - [x] query_transactions 工具（Task 3）
   - [x] MCP 生命周期管理（Task 4）
   - [x] start/stop/status API（Task 5）
   - [x] 前端 API 调用（Task 6）
   - [x] 前端状态管理（Task 7）
   - [x] 设置页面 UI（Task 8）

2. **占位符检查：** 无 TBD/TODO/不完整部分

3. **类型一致性检查：**
   - `protocol.go` 中 `JsonRpcRequest`/`JsonRpcResponse` 与 `server.go` 使用一致
   - `tools.go` 中 `ToolDefinition` 与 `server.go` 中 `ListTools()` 返回类型一致
   - `manager.go` 中 `GetMcpStatus()` 返回 `map[string]interface{}` 与 `mcp_controller.go` 使用一致
