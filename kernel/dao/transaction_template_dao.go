package dao

import (
	"sync"

	"github.com/billadm/models"
	"github.com/billadm/workspace"
)

var (
	trTemplateDao     TransactionTemplateDao
	trTemplateDaoOnce sync.Once
)

func GetTrTemplateDao() TransactionTemplateDao {
	if trTemplateDao != nil {
		return trTemplateDao
	}

	trTemplateDaoOnce.Do(func() {
		trTemplateDao = &transactionTemplateDaoImpl{}
	})

	return trTemplateDao
}

type TransactionTemplateDao interface {
	Create(ws *workspace.Workspace, template *models.TransactionTemplate) error
	DeleteById(ws *workspace.Workspace, templateId string) error
	ListByLedgerId(ws *workspace.Workspace, ledgerId string) ([]*models.TransactionTemplate, error)
	GetById(ws *workspace.Workspace, templateId string) (*models.TransactionTemplate, error)
	UpdateSortOrder(ws *workspace.Workspace, templateId string, sortOrder int) error
	GetMaxSortOrder(ws *workspace.Workspace, ledgerId string) (int, error)
}

var _ TransactionTemplateDao = &transactionTemplateDaoImpl{}

type transactionTemplateDaoImpl struct{}

func (t *transactionTemplateDaoImpl) Create(ws *workspace.Workspace, template *models.TransactionTemplate) error {
	return ws.GetDb().Create(template).Error
}

func (t *transactionTemplateDaoImpl) DeleteById(ws *workspace.Workspace, templateId string) error {
	return ws.GetDb().
		Where("template_id = ?", templateId).
		Delete(&models.TransactionTemplate{}).Error
}

func (t *transactionTemplateDaoImpl) ListByLedgerId(ws *workspace.Workspace, ledgerId string) ([]*models.TransactionTemplate, error) {
	templates := make([]*models.TransactionTemplate, 0)
	err := ws.GetDb().
		Where("ledger_id = ?", ledgerId).
		Order("sort_order ASC, created_at DESC").
		Find(&templates).Error
	return templates, err
}

func (t *transactionTemplateDaoImpl) GetById(ws *workspace.Workspace, templateId string) (*models.TransactionTemplate, error) {
	var template models.TransactionTemplate
	err := ws.GetDb().
		Where("template_id = ?", templateId).
		First(&template).Error
	return &template, err
}

func (t *transactionTemplateDaoImpl) UpdateSortOrder(ws *workspace.Workspace, templateId string, sortOrder int) error {
	return ws.GetDb().
		Model(&models.TransactionTemplate{}).
		Where("template_id = ?", templateId).
		Update("sort_order", sortOrder).Error
}

func (t *transactionTemplateDaoImpl) GetMaxSortOrder(ws *workspace.Workspace, ledgerId string) (int, error) {
	var maxSortOrder int
	err := ws.GetDb().Model(&models.TransactionTemplate{}).
		Where("ledger_id = ?", ledgerId).
		Select("COALESCE(MAX(sort_order), 0)").
		Scan(&maxSortOrder).Error
	if err != nil {
		return 0, err
	}
	return maxSortOrder, nil
}