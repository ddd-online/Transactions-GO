package dto

import "github.com/billadm/models"

type TagDto struct {
	Name     string `json:"name"`
	Scope    string `json:"scope"`
	Category string `json:"category"`
}

// ToTag 将 TagDto 转换为 Tag
// 用于将前端传入的数据保存到数据库
func (dto *TagDto) ToTag() *models.Tag {
	return &models.Tag{
		Name:     dto.Name,
		Scope:    dto.Scope,
		Category: dto.Category,
	}
}

// FromTag 从 Tag 模型填充 TagDto
// 用于将数据库数据转换为前端可用的 DTO
func (dto *TagDto) FromTag(tag *models.Tag) {
	dto.Name = tag.Name
	dto.Scope = tag.Scope
	dto.Category = tag.Category
}
