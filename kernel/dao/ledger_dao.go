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
	ModifyLedger(ws *workspace.Workspace, ledger *models.Ledger) error
	ListAllLedger(ws *workspace.Workspace) ([]models.Ledger, error)
	QueryLedgerById(ws *workspace.Workspace, ledgerId string) (*models.Ledger, error)
	DeleteLedgerById(ws *workspace.Workspace, ledgerId string) error
}

var _ LedgerDao = &ledgerDaoImpl{}

type ledgerDaoImpl struct{}

func (l *ledgerDaoImpl) CreateLedger(ws *workspace.Workspace, ledger *models.Ledger) error {
	return ws.GetDb().Create(ledger).Error
}

func (l *ledgerDaoImpl) ModifyLedger(ws *workspace.Workspace, ledger *models.Ledger) error {
	return ws.GetDb().Model(ledger).
		Where("id = ?", ledger.ID).
		Updates(map[string]interface{}{
			"name":        ledger.Name,
			"description": ledger.Description,
		}).Error
}

func (l *ledgerDaoImpl) ListAllLedger(ws *workspace.Workspace) ([]models.Ledger, error) {
	ledgers := make([]models.Ledger, 0)
	err := ws.GetDb().Find(&ledgers).Error
	return ledgers, err
}

func (l *ledgerDaoImpl) QueryLedgerById(ws *workspace.Workspace, ledgerId string) (*models.Ledger, error) {
	var ledger models.Ledger
	err := ws.GetDb().Where("id = ?", ledgerId).First(&ledger).Error
	return &ledger, err
}

func (l *ledgerDaoImpl) DeleteLedgerById(ws *workspace.Workspace, ledgerId string) error {
	return ws.GetDb().Where("id = ?", ledgerId).Delete(&models.Ledger{}).Error
}
