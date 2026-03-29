package dto

import "github.com/billadm/models"

type TagDto struct {
	Name                      string `json:"name"`
	CategoryTransactionType   string `json:"categoryTransactionType"`
}

// ToTag 将 TagDto 转换为 Tag
// 用于将前端传入的数据保存到数据库
func (dto *TagDto) ToTag() *models.Tag {
	return &models.Tag{
		Name:                      dto.Name,
		CategoryTransactionType:   dto.CategoryTransactionType,
	}
}

// FromTag 从 Tag 模型填充 TagDto
// 用于将数据库数据转换为前端可用的 DTO
func (dto *TagDto) FromTag(tag *models.Tag) {
	dto.Name = tag.Name
	dto.CategoryTransactionType = tag.CategoryTransactionType
}
