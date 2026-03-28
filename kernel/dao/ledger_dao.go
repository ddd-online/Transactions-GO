package dao

import (
	"sync"

	"github.com/billadm/models"
	"github.com/billadm/workspace"
)

var (
	ledgerDao     LedgerDao
	ledgerDaoOnce sync.Once
)

func GetLedgerDao() LedgerDao {
	if ledgerDao != nil {
		return ledgerDao
	}
	ledgerDaoOnce.Do(func() {
		ledgerDao = &ledgerDaoImpl{}
	})
	return ledgerDao
}

type LedgerDao interface {
	CreateLedger(ws *workspace.Workspace, ledger *models.Ledger) error
	ModifyLedgerName(ws *workspace.Workspace, ledger *models.Ledger) error
	ListAllLedger(ws *workspace.Workspace) ([]models.Ledger, error)
	QueryLedgerById(ws *workspace.Workspace, ledgerId string) (*models.Ledger, error)
	DeleteLedgerById(ws *workspace.Workspace, ledgerId string) error
}

var _ LedgerDao = &ledgerDaoImpl{}

type ledgerDaoImpl struct{}

func (l *ledgerDaoImpl) CreateLedger(ws *workspace.Workspace, ledger *models.Ledger) error {
	if err := ws.GetDb().Create(ledger).Error; err != nil {
		return err
	}

	return nil
}

func (l *ledgerDaoImpl) ModifyLedgerName(ws *workspace.Workspace, ledger *models.Ledger) error {
	if err := ws.GetDb().Model(ledger).
		Where("id = ?", ledger.ID).
		Update("name", ledger.Name).Error; err != nil {
		return err
	}

	return nil
}

func (l *ledgerDaoImpl) ListAllLedger(ws *workspace.Workspace) ([]models.Ledger, error) {
	ledgers := make([]models.Ledger, 0)
	if err := ws.GetDb().Find(&ledgers).Error; err != nil {
		return nil, err
	}

	return ledgers, nil
}

func (l *ledgerDaoImpl) QueryLedgerById(ws *workspace.Workspace, ledgerId string) (*models.Ledger, error) {
	var ledger models.Ledger
	if err := ws.GetDb().Where("id = ?", ledgerId).First(&ledger).Error; err != nil {
		return nil, err
	}

	return &ledger, nil
}

func (l *ledgerDaoImpl) DeleteLedgerById(ws *workspace.Workspace, ledgerId string) error {
	if err := ws.GetDb().Where("id = ?", ledgerId).Delete(&models.Ledger{}).Error; err != nil {
		return err
	}

	return nil
}
