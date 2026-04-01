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
}

var _ TransactionTemplateDao = &transactionTemplateDaoImpl{}

type transactionTemplateDaoImpl struct{}

func (t *transactionTemplateDaoImpl) Create(ws *workspace.Workspace, template *models.TransactionTemplate) error {
	if err := ws.GetDb().Create(template).Error; err != nil {
		return err
	}
	return nil
}

func (t *transactionTemplateDaoImpl) DeleteById(ws *workspace.Workspace, templateId string) error {
	if err := ws.GetDb().
		Where("template_id = ?", templateId).
		Delete(&models.TransactionTemplate{}).Error; err != nil {
		return err
	}
	return nil
}

func (t *transactionTemplateDaoImpl) ListByLedgerId(ws *workspace.Workspace, ledgerId string) ([]*models.TransactionTemplate, error) {
	templates := make([]*models.TransactionTemplate, 0)
	if err := ws.GetDb().
		Where("ledger_id = ?", ledgerId).
		Order("created_at desc").
		Find(&templates).Error; err != nil {
		return nil, err
	}
	return templates, nil
}

func (t *transactionTemplateDaoImpl) GetById(ws *workspace.Workspace, templateId string) (*models.TransactionTemplate, error) {
	var template models.TransactionTemplate
	if err := ws.GetDb().
		Where("template_id = ?", templateId).
		First(&template).Error; err != nil {
		return nil, err
	}
	return &template, nil
}