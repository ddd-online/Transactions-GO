package mcp

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/billadm/models/dto"
	"github.com/billadm/service"
	"github.com/billadm/workspace"
	"github.com/sirupsen/logrus"
)

// ToolDefinition defines a tool exposed by the server
type ToolDefinition struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	InputSchema map[string]interface{} `json:"inputSchema"`
}

// ToolHandler handles tool-related operations
type ToolHandler struct {
	tools []ToolDefinition
}

// NewToolHandler creates a new ToolHandler instance
func NewToolHandler() *ToolHandler {
	h := &ToolHandler{
		tools: make([]ToolDefinition, 0, 2),
	}

	// query_ledgers: Query all ledgers
	h.tools = append(h.tools, ToolDefinition{
		Name:        "query_ledgers",
		Description: "Query all ledgers in the current workspace. Returns a list of all ledgers with their IDs, names, and descriptions.",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{},
		},
	})

	// query_transactions: Query transactions with filters
	h.tools = append(h.tools, ToolDefinition{
		Name:        "query_transactions",
		Description: "Query transaction records with optional filters. Supports filtering by ledger, time range, transaction type, category, tags, and description. Supports pagination and sorting.",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"ledger_id": map[string]interface{}{
					"type":        "string",
					"description": "The ledger ID to query transactions from (required)",
				},
				"time_range": map[string]interface{}{
					"type":        "array",
					"items":       map[string]interface{}{"type": "number"},
					"description": "Time range as [start_timestamp, end_timestamp] in milliseconds",
				},
				"transaction_type": map[string]interface{}{
					"type":        "string",
					"description": "Transaction type filter: 'expense', 'income', or 'transfer'",
				},
				"category": map[string]interface{}{
					"type":        "string",
					"description": "Category name to filter by",
				},
				"tags": map[string]interface{}{
					"type":        "array",
					"items":       map[string]interface{}{"type": "string"},
					"description": "Tags to filter by",
				},
				"description": map[string]interface{}{
					"type":        "string",
					"description": "Description substring to search for",
				},
				"offset": map[string]interface{}{
					"type":        "integer",
					"description": "Pagination offset (default: 0)",
				},
				"limit": map[string]interface{}{
					"type":        "integer",
					"description": "Pagination limit (default: 20, max: 100)",
				},
				"sort_fields": map[string]interface{}{
					"type":        "array",
					"items":       map[string]interface{}{"type": "object"},
					"description": "Sort fields, e.g. [{\"field\": \"transactionAt\", \"order\": \"desc\"}]",
				},
			},
			"required": []string{"ledger_id"},
		},
	})

	return h
}

// ListTools returns the list of available tools
func (h *ToolHandler) ListTools() []ToolDefinition {
	return h.tools
}

// CallTool calls a tool by name with the given arguments
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

// queryLedgers queries all ledgers in the current workspace
func (h *ToolHandler) queryLedgers() (string, error) {
	ws := workspace.Manager.OpenedWorkspace()
	if ws == nil {
		return "", fmt.Errorf("工作空间未打开")
	}

	ledgers, err := service.GetLedgerService().ListAllLedger(ws)
	if err != nil {
		logrus.Errorf("query ledgers failed: %v", err)
		return "", fmt.Errorf("查询账本失败: %v", err)
	}

	if len(ledgers) == 0 {
		return "未找到账本", nil
	}

	var sb strings.Builder
	for _, ledger := range ledgers {
		sb.WriteString(fmt.Sprintf("[%s] %s", ledger.ID, ledger.Name))
		if ledger.Description != "" {
			sb.WriteString(fmt.Sprintf(" - %s", ledger.Description))
		}
		sb.WriteString("\n")
	}

	return sb.String(), nil
}

// queryTransactionsArgs defines the arguments for query_transactions tool
type queryTransactionsArgs struct {
	LedgerID         string   `json:"ledger_id"`
	TimeRange        []int64  `json:"time_range,omitempty"`
	TransactionType  string   `json:"transaction_type,omitempty"`
	Category         string   `json:"category,omitempty"`
	Tags             []string `json:"tags,omitempty"`
	Description      string   `json:"description,omitempty"`
	Offset           int      `json:"offset,omitempty"`
	Limit            int      `json:"limit,omitempty"`
	SortFields       []struct {
		Field string `json:"field"`
		Order string `json:"order"`
	} `json:"sort_fields,omitempty"`
}

// queryTransactions queries transaction records with filters
func (h *ToolHandler) queryTransactions(arguments json.RawMessage) (string, error) {
	ws := workspace.Manager.OpenedWorkspace()
	if ws == nil {
		return "", fmt.Errorf("工作空间未打开")
	}

	var args queryTransactionsArgs
	if err := json.Unmarshal(arguments, &args); err != nil {
		return "", fmt.Errorf("invalid arguments: %v", err)
	}

	// Build query condition
	condition := &dto.TrQueryCondition{
		LedgerID: args.LedgerID,
		Offset:   args.Offset,
		Limit:    args.Limit,
		TsRange:  args.TimeRange,
		Items:    make([]dto.QueryConditionItem, 0),
	}

	// Apply default pagination if not specified
	if condition.Offset < 0 {
		condition.Offset = 0
	}
	if condition.Limit <= 0 {
		condition.Limit = 20
	}
	if condition.Limit > 100 {
		condition.Limit = 100
	}

	// Build condition items
	if args.TransactionType != "" || args.Category != "" || len(args.Tags) > 0 || args.Description != "" {
		item := dto.QueryConditionItem{
			TransactionType: args.TransactionType,
			Category:        args.Category,
			Tags:            args.Tags,
			Description:     args.Description,
		}
		condition.Items = append(condition.Items, item)
	}

	// Convert sort fields
	if len(args.SortFields) > 0 {
		condition.SortFields = make([]dto.QueryConditionSortField, 0, len(args.SortFields))
		for _, sf := range args.SortFields {
			condition.SortFields = append(condition.SortFields, dto.QueryConditionSortField{
				Field: sf.Field,
				Order: sf.Order,
			})
		}
	}

	result, err := service.GetTrService().QueryTrsOnCondition(ws, condition)
	if err != nil {
		logrus.Errorf("query transactions failed: %v", err)
		return "", fmt.Errorf("查询交易记录失败: %v", err)
	}

	// Format results
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("共 %d 条记录 (第 %d-%d 项)\n\n",
		result.Total, condition.Offset+1, condition.Offset+len(result.Items)))

	if len(result.Items) == 0 {
		return "未找到交易记录", nil
	}

	for _, tr := range result.Items {
		sb.WriteString(formatTransactionRecord(tr))
		sb.WriteString("\n")
	}

	return sb.String(), nil
}

// formatTransactionRecord formats a single transaction record for display
func formatTransactionRecord(tr *dto.TransactionRecordDto) string {
	var sb strings.Builder

	// Format timestamp
	transactionAt := time.Unix(tr.TransactionAt/1000, 0).Format("2006-01-02 15:04:05")

	sb.WriteString(fmt.Sprintf("[%s] %s | %s | %d",
		tr.TransactionID[:8],
		transactionAt,
		tr.TransactionType,
		tr.Price))

	if tr.Category != "" {
		sb.WriteString(fmt.Sprintf(" | %s", tr.Category))
	}

	if len(tr.Tags) > 0 {
		sb.WriteString(fmt.Sprintf(" | 标签: %s", strings.Join(tr.Tags, ", ")))
	}

	if tr.Description != "" {
		sb.WriteString(fmt.Sprintf(" | %s", tr.Description))
	}

	return sb.String()
}

