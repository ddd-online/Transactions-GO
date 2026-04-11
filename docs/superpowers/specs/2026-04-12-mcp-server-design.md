# MCP Server 设计文档

## 概述

为 Billadm-GO 增加 MCP（Model Context Protocol）服务器支持，使 AI 助手（如 Claude Code）能够通过 TCP 连接查询账本和交易记录数据。

## 架构

```
┌─────────────────┐      TCP (127.0.0.1:31944)      ┌─────────────────┐
│   Claude Code   │◄───────────────────────────────►│   Go MCP Server │
│   (MCP Client)  │                                 │   (goroutine)    │
└─────────────────┘                                 └────────┬────────┘
                                                              │
                                          ┌───────────────────┼───────────────────┐
                                          │                   │                   │
                                          ▼                   ▼                   ▼
                                   ┌────────────┐    ┌──────────────┐    ┌──────────────┐
                                   │ Ledger Svc │    │ TrRecord Svc │    │  Workspace   │
                                   └────────────┘    └──────────────┘    └──────────────┘
```

## MCP 协议实现

### 协议基础

- 基于 JSON-RPC 2.0
- 通过 TCP Socket 通信，监听 `127.0.0.1:31944`
- 每个连接一个 goroutine 处理
- 使用 `bufio.NewReader` + `json.Decoder` 解析请求
- 使用 `json.NewEncoder` 发送响应

### 核心方法

#### initialize

Claude Code 连接时首先调用，服务器返回协议版本和能力。

```json
// Request
{"jsonrpc": "2.0", "id": 1, "method": "initialize", "params": {"protocolVersion": "2024-11-05", "capabilities": {}}}

// Response
{"jsonrpc": "2.0", "id": 1, "result": {"protocolVersion": "2024-11-05", "capabilities": {"tools": {}}}}
```

#### tools/list

返回所有可用工具。

```json
// Request
{"jsonrpc": "2.0", "id": 2, "method": "tools/list", "params": {}}

// Response
{"jsonrpc": "2.0", "id": 2, "result": {"tools": [
  {"name": "query_ledgers", "description": "查询所有账本", "inputSchema": {"type": "object", "properties": {}}},
  {"name": "query_transactions", "description": "查询交易记录", "inputSchema": {"type": "object", "properties": {
    "ledger_id": {"type": "string", "description": "账本ID（必填）"},
    "time_range": {"type": "array", "items": {"type": "integer"}, "description": "时间戳范围 [start, end]"},
    "transaction_type": {"type": "string", "enum": ["expense", "income", "transfer"]},
    "category": {"type": "string"},
    "tags": {"type": "array", "items": {"type": "string"}},
    "description": {"type": "string"},
    "offset": {"type": "integer", "default": 0},
    "limit": {"type": "integer", "default": 50},
    "sort_fields": {"type": "array", "items": {"type": "object", "properties": {"field": {"type": "string"}, "order": {"type": "string", "enum": ["asc", "desc"]}}}}
  }, "required": ["ledger_id"]}}
]}}
```

#### tools/call

调用具体工具执行查询。

```json
// Request
{"jsonrpc": "2.0", "id": 3, "method": "tools/call", "params": {"name": "query_ledgers", "arguments": {}}}

// Response
{"jsonrpc": "2.0", "id": 3, "result": {"content": [{"type": "text", "text": "账本列表：\n1. 日常账本 (id: xxx, 描述: xxx)\n2. ..."}]}}
```

## MCP 工具实现

### query_ledgers

查询所有账本，无参数。

### query_transactions

按条件查询交易记录。

**参数：**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| ledger_id | string | 是 | 账本ID |
| time_range | [int64, int64] | 否 | 时间戳范围 |
| transaction_type | string | 否 | expense/income/transfer |
| category | string | 否 | 分类名称 |
| tags | string[] | 否 | 标签列表 |
| description | string | 否 | 备注包含的字符 |
| offset | int | 否 | 分页偏移，默认0 |
| limit | int | 否 | 每页数量，默认50 |
| sort_fields | {field, order}[] | 否 | 排序字段 |

**返回格式：**

```
交易记录列表：
1. [2024-01-15] 支出 - 餐饮 - 午餐 - #工作 #餐饮 - 25.50元 - 备注: 午餐
2. [2024-01-16] 收入 - 工资 - 1月工资 - #工资 - 5000.00元 - 备注: 
...
共 125 条记录，第 1/3 页
```

## 前端控制 API

### 接口设计

```
POST /api/v1/mcp/start  - 启动 MCP 服务器
POST /api/v1/mcp/stop   - 停止 MCP 服务器
GET  /api/v1/mcp/status - 获取 MCP 服务器状态
```

### 响应格式

```json
{
  "code": 0,
  "msg": "",
  "data": {
    "running": true,
    "port": 31944
  }
}
```

### 错误码

| code | 说明 |
|------|------|
| 0 | 成功 |
| -1 | 失败（已在运行/未运行等） |

## 前端实现

### 状态管理 (mcpStore.ts)

```typescript
interface McpState {
  isRunning: boolean;
  port: number;
}

interface McpStore {
  state: McpState;
  start(): Promise<void>;
  stop(): Promise<void>;
  refreshStatus(): Promise<void>;
}
```

### 设置页面

在现有设置页面增加 MCP 区块：

```
┌─────────────────────────────────┐
│ MCP 服务                    [关] │
│ 启用后 AI 助手可查询账本和交易记录 │
│ 端口: 31944                     │
└─────────────────────────────────┘
```

- 开关按钮控制 MCP 服务的启停
- 显示当前服务状态
- 显示服务端口

## 文件结构

```
kernel/
  mcp/
    server.go       # MCP TCP 服务器，JSON-RPC 处理
    tools.go        # MCP 工具实现（账本/交易查询）
    protocol.go     # JSON-RPC 2.0 协议工具函数

app/src/
  stores/
    mcpStore.ts     # MCP 状态管理
  backend/
    mcp.ts          # MCP 控制 API 调用

app/src/components/  (现有设置页面中增加 MCP 区块)
```

## 实现步骤

1. 创建 `kernel/mcp/` 目录，实现 MCP 协议处理
2. 实现 `tools/list` 和 `tools/call` 方法
3. 实现 `query_ledgers` 和 `query_transactions` 工具
4. 在现有 HTTP API 中增加 MCP 控制接口
5. 前端增加 mcpStore 和设置页面 UI
6. 集成测试

## 约束与限制

- 仅支持查询，不支持创建/修改/删除
- MCP 服务器只能有一个运行实例
- MCP 服务端口固定为 `31944`
- 需要先打开工作空间才能启动 MCP 服务
