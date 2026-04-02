package dto

import "github.com/billadm/models"

type CategoryDto struct {
	Name            string `json:"name"`
	TransactionType string `json:"transactionType"`
	SortOrder       int    `json:"sortOrder"`
}

type CreateCategoryRequest struct {
	LedgerID        string `json:"ledgerId"`
	Name            string `json:"name"`
	TransactionType string `json:"transactionType"`
	SortOrder       int    `json:"sortOrder"`
}

type UpdateCategorySortRequest struct {
	LedgerID        string `json:"ledgerId"`
	Name            string `json:"name"`
	TransactionType string `json:"transactionType"`
	SortOrder       int    `json:"sortOrder"`
}

// ToCategory 将 CategoryDto 转换为 Category
func (dto *CategoryDto) ToCategory() *models.Category {
	return &models.Category{
		Name:            dto.Name,
		TransactionType: dto.TransactionType,
		SortOrder:       dto.SortOrder,
	}
}

// FromCategory 从 Category 填充 CategoryDto
func (dto *CategoryDto) FromCategory(category *models.Category) {
	dto.Name = category.Name
	dto.TransactionType = category.TransactionType
	dto.SortOrder = category.SortOrder
}