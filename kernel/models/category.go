package models

type Category struct {
	Name            string `gorm:"primaryKey;comment:消费类型" json:"name"`
	Scope           string `gorm:"not null;comment:作用域" json:"scope"`
	TransactionType string `gorm:"not null;comment:交易类型" json:"transaction_type"`
}

func (tr *Category) TableName() string {
	return "tbl_billadm_category"
}
