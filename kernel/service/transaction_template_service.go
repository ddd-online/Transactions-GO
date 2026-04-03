package service

import (
	"sync"

	"github.com/billadm/dao"
	"github.com/billadm/models/dto"
	"github.com/billadm/util"
	"github.com/billadm/workspace"
	"github.com/sirupsen/logrus"
)

var (
	trTemplateService     TransactionTemplateService
	trTemplateServiceOnce sync.Once
)

func GetTrTemplateService() TransactionTemplateService {
	if trTemplateService != nil {
		return trTemplateService
	}

	trTemplateServiceOnce.Do(func() {
		trTemplateService = &transactionTemplateServiceImpl{
			trTemplateDao: dao.GetTrTemplateDao(),
		}
	})

	return trTemplateService
}

type TransactionTemplateService interface {
	Create(ws *workspace.Workspace, dto *dto.TransactionTemplateDto) (string, error)
	DeleteById(ws *workspace.Workspace, templateId string) error
	ListByLedgerId(ws *workspace.Workspace, ledgerId string) ([]*dto.TransactionTemplateDto, error)
	UpdateSortOrder(ws *workspace.Workspace, templateId string, ledgerId string, sortOrder int) error
}

var _ TransactionTemplateService = &transactionTemplateServiceImpl{}

type transactionTemplateServiceImpl struct {
	trTemplateDao dao.TransactionTemplateDao
}

func (t *transactionTemplateServiceImpl) Create(ws *workspace.Workspace, dto *dto.TransactionTemplateDto) (string, error) {
	logrus.Infof("start to create transaction template, ledger id: %s, name: %s", dto.LedgerID, dto.TemplateName)

	templateID := util.GetUUID()

	// Get max sort order for this ledger
	maxSortOrder, err := t.trTemplateDao.GetMaxSortOrder(ws, dto.LedgerID)
	if err != nil {
		logrus.Errorf("get max sort order failed: %v", err)
		return "", err
	}

	record := dto.ToTransactionTemplate()
	record.TemplateID = templateID
	record.SortOrder = maxSortOrder + 1

	if err := t.trTemplateDao.Create(ws, record); err != nil {
		logrus.Errorf("create transaction template failed: %v", err)
		return "", err
	}

	logrus.Infof("create transaction template success, ledger id: %s, name: %s", dto.LedgerID, dto.TemplateName)
	return templateID, nil
}

func (t *transactionTemplateServiceImpl) DeleteById(ws *workspace.Workspace, templateId string) error {
	logrus.Infof("start to delete transaction template, id: %s", templateId)

	if err := t.trTemplateDao.DeleteById(ws, templateId); err != nil {
		logrus.Errorf("delete transaction template failed: %v", err)
		return err
	}

	logrus.Infof("delete transaction template success, id: %s", templateId)
	return nil
}

func (t *transactionTemplateServiceImpl) ListByLedgerId(ws *workspace.Workspace, ledgerId string) ([]*dto.TransactionTemplateDto, error) {
	logrus.Infof("start to list transaction templates, ledger id: %s", ledgerId)

	templates, err := t.trTemplateDao.ListByLedgerId(ws, ledgerId)
	if err != nil {
		return nil, err
	}

	dtos := make([]*dto.TransactionTemplateDto, 0, len(templates))
	for _, template := range templates {
		dto := &dto.TransactionTemplateDto{}
		dto.FromTransactionTemplate(template)
		dtos = append(dtos, dto)
	}

	logrus.Infof("list transaction templates success, ledger id: %s, count: %d", ledgerId, len(dtos))
	return dtos, nil
}

func (t *transactionTemplateServiceImpl) UpdateSortOrder(ws *workspace.Workspace, templateId string, ledgerId string, sortOrder int) error {
	logrus.Infof("start to update template sort, templateId: %s, sortOrder: %d", templateId, sortOrder)

	// Reindex all templates for this ledger to ensure sequential sort_order values
	templates, err := t.trTemplateDao.ListByLedgerId(ws, ledgerId)
	if err != nil {
		logrus.Errorf("list templates failed: %v", err)
		return err
	}

	// Reassign sort_order from 0 based on current order
	for i, template := range templates {
		if template.SortOrder != i {
			if err := t.trTemplateDao.UpdateSortOrder(ws, template.TemplateID, i); err != nil {
				logrus.Errorf("reindex template sort failed: %v", err)
				return err
			}
		}
	}

	if err := t.trTemplateDao.UpdateSortOrder(ws, templateId, sortOrder); err != nil {
		logrus.Errorf("update template sort failed: %v", err)
		return err
	}

	logrus.Infof("update template sort success, templateId: %s", templateId)
	return nil
}