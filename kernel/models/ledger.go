package models

type Ledger struct {
	ID        string `gorm:"primaryKey;comment:账本UUID" json:"id"`
	Name      string `gorm:"not null;comment:账本名称" json:"name"`
	CreatedAt int64  `gorm:"autoCreateTime:unix;not null;comment:创建时间" json:"created_at"`
	UpdatedAt int64  `gorm:"autoUpdateTime:unix;not null;comment:更新时间" json:"updated_at"`
}

func (l *Ledger) TableName() string {
	return "tbl_billadm_ledger"
}
