package models

// ChartLine 图表曲线配置（与前端 ChartLine 对应）
type ChartLine struct {
	Label          string                  `json:"label"`
	TransactionType string                `json:"transactionType"`
	IncludeOutlier bool                   `json:"includeOutlier"`
	Conditions     []QueryConditionItem   `json:"conditions"`
}

// QueryConditionItem 查询条件项
type QueryConditionItem struct {
	TransactionType string   `json:"transactionType"`
	Category       string   `json:"category"`
	Tags           []string `json:"tags"`
	TagPolicy      string   `json:"tagPolicy"`
	TagNot         bool     `json:"tagNot"`
	Description    string   `json:"description"`
}

// Chart 图表配置
type Chart struct {
	ChartID      string `gorm:"primaryKey;comment:图表UUID" json:"chart_id"`
	LedgerID     string `gorm:"not null;comment:关联账本ID" json:"ledger_id"`
	Title        string `gorm:"not null;comment:图表名称" json:"title"`
	Granularity  string `gorm:"not null;comment:时间粒度 year/month" json:"granularity"`
	ChartLines   string `gorm:"not null;comment:曲线配置JSON" json:"chart_lines"`
	ChartType    string `gorm:"not null;default:'line';comment:图表类型 line/bar" json:"chart_type"`
	IsPreset     bool   `gorm:"not null;default:false;comment:是否预设图表" json:"is_preset"`
	SortOrder    int    `gorm:"not null;default:0;comment:排序" json:"sort_order"`
	CreatedAt    int64  `gorm:"autoCreateTime:unix;not null;comment:创建时间" json:"created_at"`
	UpdatedAt    int64  `gorm:"autoUpdateTime:unix;not null;comment:更新时间" json:"updated_at"`
}

func (c *Chart) TableName() string {
	return "tbl_billadm_chart"
}
