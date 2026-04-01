package models

// TransactionTemplate 消费记录模板
type TransactionTemplate struct {
	TemplateID      string `gorm:"primaryKey;comment:模板UUID" json:"template_id"`
	LedgerID        string `gorm:"not null;comment:关联账本ID" json:"ledger_id"`
	TemplateName    string `gorm:"not null;comment:模板名称" json:"template_name"`
	TransactionType string `gorm:"not null;comment:交易类型" json:"transaction_type"`
	Category        string `gorm:"not null;comment:分类ID" json:"category"`
	Tags            string `gorm:"comment:标签集，JSON数组格式" json:"tags"`
	Flags           string `gorm:"comment:标记集" json:"flags"`
	Description     string `gorm:"comment:交易描述" json:"description"`
	CreatedAt       int64  `gorm:"autoCreateTime:unix;not null;comment:创建时间" json:"created_at"`
	UpdatedAt       int64  `gorm:"autoUpdateTime:unix;not null;comment:更新时间" json:"updated_at"`
}

func (tt *TransactionTemplate) TableName() string {
	return "tbl_billadm_transaction_tpl"
}