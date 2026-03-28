package dto

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/billadm/models"
)

const (
	Any = "any"
	All = "all"
	Not = "not"
)

func JsonQueryCondition(c *gin.Context, result *models.Result) (*TrQueryCondition, bool) {
	ret := &TrQueryCondition{
		Offset:  -1,
		Limit:   -1,
		TsRange: make([]int64, 0),
		Items:   make([]QueryConditionItem, 0),
	}
	if err := c.BindJSON(ret); nil != err {
		result.Code = -1
		result.Msg = fmt.Sprintf("解析消费记录查询条件失败: %v", err)
		return nil, false
	}
	return ret, true
}

type TrQueryCondition struct {
	LedgerID string               `json:"ledgerId"`
	Offset   int                  `json:"offset"`
	Limit    int                  `json:"limit"`
	TsRange  []int64              `json:"tsRange"`
	Items    []QueryConditionItem `json:"items"`
}

func (qc *TrQueryCondition) Validate(result *models.Result) bool {
	if len(qc.LedgerID) == 0 {
		result.Code = -1
		result.Msg = fmt.Sprintf("账本Id不可为空: %s", qc.LedgerID)
		return false
	}
	return true
}

type QueryConditionItem struct {
	TransactionType string   `json:"transactionType"`
	Category        string   `json:"category"`
	Tags            []string `json:"tags"`
	TagPolicy       string   `json:"tagPolicy"`   // 如何匹配tag列表
	TagNot          bool     `json:"tagNot"`      // 是否对tag匹配策略取反
	Description     string   `json:"description"` // 描述包含的字符
}
