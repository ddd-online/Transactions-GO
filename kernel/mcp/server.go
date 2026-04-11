package mcp

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/billadm/models/dto"
	"github.com/billadm/service"
	"github.com/billadm/workspace"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/sirupsen/logrus"
)

const (
	ServerPort = 31944
)

// McpServer is the main MCP server using mcp-go SDK
type McpServer struct {
	httpServer *server.StreamableHTTPServer
	mcpServer  *server.MCPServer
	mu          sync.Mutex
	running     bool
}

// NewMcpServer creates a new MCP server instance using mcp-go SDK
func NewMcpServer() *McpServer {
	s := server.NewMCPServer("Billadm MCP", "1.0.0",
		server.WithToolCapabilities(true),
	)

	// Add query_ledgers tool
	s.AddTool(
		mcp.NewTool("query_ledgers",
			mcp.WithDescription("查询所有账本"),
		),
		queryLedgersHandler,
	)

	// Add query_transactions tool
	s.AddTool(
		mcp.NewTool("query_transactions",
			mcp.WithDescription("按条件查询交易记录"),
			mcp.WithString("ledger_id", mcp.Required(), mcp.Description("账本ID")),
			mcp.WithArray("time_range",
				mcp.Description("时间戳范围 [start, end]"),
				mcp.WithNumberItems(),
			),
			mcp.WithString("transaction_type", mcp.Description("expense/income/transfer")),
			mcp.WithString("category", mcp.Description("分类名称")),
			mcp.WithArray("tags",
				mcp.Description("标签列表"),
				mcp.WithStringItems(),
			),
			mcp.WithString("description", mcp.Description("备注包含的字符")),
			mcp.WithNumber("offset", mcp.Description("分页偏移")),
			mcp.WithNumber("limit", mcp.Description("每页数量，最大100")),
		),
		queryTransactionsHandler,
	)

	return &McpServer{
		mcpServer: s,
	}
}

// Start starts the MCP server in a goroutine
func (s *McpServer) Start() error {
	s.mu.Lock()
	if s.running {
		s.mu.Unlock()
		return nil
	}

	httpServer := server.NewStreamableHTTPServer(s.mcpServer)

	mux := http.NewServeMux()
	mux.Handle("/mcp", httpServer)

	addr := fmt.Sprintf("127.0.0.1:%d", ServerPort)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		s.mu.Unlock()
		return fmt.Errorf("failed to listen on %s: %w", addr, err)
	}

	// Start HTTP server in a goroutine so it doesn't block
	go func() {
		if err := http.Serve(listener, mux); err != nil {
			logrus.Errorf("MCP HTTP server error: %v", err)
		}
	}()

	s.httpServer = httpServer
	s.running = true
	s.mu.Unlock()

	logrus.Infof("MCP server listening on %s/mcp", addr)
	return nil
}

// Stop stops the server
func (s *McpServer) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.running {
		return nil
	}
	s.running = false
	logrus.Info("MCP server stopped")
	return nil
}

// IsRunning returns the running status of the server
func (s *McpServer) IsRunning() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.running
}

// Tool handlers

func queryLedgersHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ws := workspace.Manager.OpenedWorkspace()
	if ws == nil {
		return mcp.NewToolResultError("工作空间未打开"), nil
	}

	ledgers, err := service.GetLedgerService().ListAllLedger(ws)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	var sb strings.Builder
	for _, ledger := range ledgers {
		sb.WriteString(fmt.Sprintf("[%s] %s", ledger.ID, ledger.Name))
		if ledger.Description != "" {
			sb.WriteString(fmt.Sprintf(" - %s", ledger.Description))
		}
		sb.WriteString("\n")
	}

	return mcp.NewToolResultText(sb.String()), nil
}

func queryTransactionsHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ws := workspace.Manager.OpenedWorkspace()
	if ws == nil {
		return mcp.NewToolResultError("工作空间未打开"), nil
	}

	ledgerID, err := request.RequireString("ledger_id")
	if err != nil {
		return mcp.NewToolResultError("ledger_id is required"), nil
	}

	condition := &dto.TrQueryCondition{
		LedgerID: ledgerID,
		Offset:   0,
		Limit:    20,
	}

	// Parse optional parameters

	// time_range - array of numbers [start, end]
	timeRange := request.GetFloatSlice("time_range", nil)
	if len(timeRange) == 2 {
		condition.TsRange = []int64{int64(timeRange[0]), int64(timeRange[1])}
	}

	// transaction_type
	transactionType := request.GetString("transaction_type", "")
	if transactionType != "" {
		condition.Items = append(condition.Items, dto.QueryConditionItem{TransactionType: transactionType})
	}

	// category
	category := request.GetString("category", "")
	if category != "" {
		condition.Items = append(condition.Items, dto.QueryConditionItem{Category: category})
	}

	// description
	description := request.GetString("description", "")
	if description != "" {
		condition.Items = append(condition.Items, dto.QueryConditionItem{Description: description})
	}

	// tags - array of strings
	tags := request.GetStringSlice("tags", nil)
	if len(tags) > 0 {
		condition.Items = append(condition.Items, dto.QueryConditionItem{Tags: tags})
	}

	// offset
	offset := request.GetInt("offset", 0)
	if offset > 0 {
		condition.Offset = offset
	}

	// limit
	limit := request.GetInt("limit", 20)
	if limit > 100 {
		limit = 100
	}
	condition.Limit = limit

	result, err := service.GetTrService().QueryTrsOnCondition(ws, condition)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("共 %d 条记录\n\n", result.Total))
	for _, tr := range result.Items {
		sb.WriteString(formatTransactionRecord(tr))
		sb.WriteString("\n")
	}

	return mcp.NewToolResultText(sb.String()), nil
}

func formatTransactionRecord(tr *dto.TransactionRecordDto) string {
	return fmt.Sprintf("[%s] %s | %s | %d",
		tr.TransactionID[:8],
		time.Unix(tr.TransactionAt/1000, 0).Format("2006-01-02 15:04:05"),
		tr.TransactionType,
		tr.Price)
}