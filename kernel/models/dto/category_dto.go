package dto

import "github.com/billadm/models"

type CategoryDto struct {
	Name            string `json:"name"`
	Scope           string `json:"scope"`
	TransactionType string `json:"transactionType"`
}

// ToCategory 将 CategoryDto 转换为 Category
func (dto *CategoryDto) ToCategory() *models.Category {
	return &models.Category{
		Name:            dto.Name,
		Scope:           dto.Scope,
		TransactionType: dto.TransactionType,
	}
}

// FromCategory 从 Category 填充 CategoryDto
func (dto *CategoryDto) FromCategory(category *models.Category) {
	dto.Name = category.Name
	dto.Scope = category.Scope
	dto.TransactionType = category.TransactionType
}
