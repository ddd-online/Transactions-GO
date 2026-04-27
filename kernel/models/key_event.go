package models

type KeyEvent struct {
	ID        string `gorm:"primaryKey;comment:事件UUID" json:"id"`
	Date      string `gorm:"uniqueIndex;comment:日期 YYYY-MM-DD" json:"date"`
	Content   string `gorm:"type:text;comment:事件内容" json:"content"`
	CreatedAt int64  `gorm:"autoCreateTime:unix;not null;comment:创建时间" json:"createdAt"`
	UpdatedAt int64  `gorm:"autoUpdateTime:unix;not null;comment:更新时间" json:"updatedAt"`
}

func (k *KeyEvent) TableName() string {
	return "tbl_billadm_key_event"
}
