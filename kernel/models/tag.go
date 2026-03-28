package models

type Tag struct {
	Name     string `gorm:"not null;uniqueIndex:idx_tag_name_scope_category;comment:标签名称" json:"name"`
	Scope    string `gorm:"not null;uniqueIndex:idx_tag_name_scope_category;comment:作用域" json:"scope"`
	Category string `gorm:"not null;uniqueIndex:idx_tag_name_scope_category;comment:分类ID" json:"category"`
}

func (tr *Tag) TableName() string {
	return "tbl_billadm_tag"
}
