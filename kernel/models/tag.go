package models

type Tag struct {
	Name                      string `gorm:"not null;uniqueIndex:idx_tag_category_transaction_type;comment:标签名称" json:"name"`
	CategoryTransactionType    string `gorm:"not null;uniqueIndex:idx_tag_category_transaction_type;comment:分类:交易类型" json:"category_transaction_type"`
}

func (tr *Tag) TableName() string {
	return "tbl_billadm_tag"
}
