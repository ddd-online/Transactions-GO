package dto

import "github.com/billadm/models"

type LedgerDto struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt int64  `json:"createdAt"`
	UpdatedAt int64  `json:"updatedAt"`
}

// ToLedger 将 LedgerDto 转换为 Ledger
// 用于将前端传入的数据转换为数据库模型
func (dto *LedgerDto) ToLedger() *models.Ledger {
	return &models.Ledger{
		ID:        dto.ID,
		Name:      dto.Name,
		CreatedAt: dto.CreatedAt,
		UpdatedAt: dto.UpdatedAt,
	}
}

// FromLedger 从 Ledger 模型填充 LedgerDto
// 用于将数据库数据转换为前端可用的 DTO
func (dto *LedgerDto) FromLedger(ledger *models.Ledger) {
	dto.ID = ledger.ID
	dto.Name = ledger.Name
	dto.CreatedAt = ledger.CreatedAt
	dto.UpdatedAt = ledger.UpdatedAt
}
