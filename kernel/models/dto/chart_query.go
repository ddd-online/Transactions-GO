package dto

import (
	"github.com/gin-gonic/gin"

	"github.com/billadm/models"
)

type ChartLineCondition struct {
	Label          string               `json:"label"`
	TransactionType string             `json:"transactionType"`
	IncludeOutlier bool                `json:"includeOutlier"`
	Conditions     []QueryConditionItem `json:"conditions"`
}

type ChartQueryRequest struct {
	LedgerID    string               `json:"ledgerId"`
	TsRange     []int64            `json:"tsRange"`
	Granularity string               `json:"granularity"` // "year" or "month"
	Lines       []ChartLineCondition `json:"lines"`
}

// ChartLineData contains filtered transaction records for a single line
type ChartLineData struct {
	Label string                     `json:"label"`
	Type  string                     `json:"type"`
	Items []*TransactionRecordDto    `json:"items"`
}

type ChartQueryResponse struct {
	Lines []ChartLineData `json:"lines"`
}

func JsonChartQuery(c *gin.Context, result *models.Result) (*ChartQueryRequest, bool) {
	ret := &ChartQueryRequest{}
	if err := c.BindJSON(ret); err != nil {
		result.Code = -1
		result.Msg = "parse chart query request failed: " + err.Error()
		return nil, false
	}
	return ret, true
}
