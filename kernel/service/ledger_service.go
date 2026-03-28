package service

import (
	"fmt"
	"sync"

	"github.com/billadm/dao"
	"github.com/billadm/models"
	"github.com/billadm/util"
	"github.com/billadm/workspace"
)

var (
	ledgerService     LedgerService
	ledgerServiceOnce sync.Once
)

func GetLedgerService() LedgerService {
	if ledgerService != nil {
		return ledgerService
	}

	ledgerServiceOnce.Do(func() {
		ledgerService = &ledgerServiceImpl{
			ledgerDao: dao.GetLedgerDao(),
			trDao:     dao.GetTrDao(),
			trTagDao:  dao.GetTrTagDao(),
		}
	})

	return ledgerService
}

type LedgerService interface {
	CreateLedger(ws *workspace.Workspace, ledgerName string) (string, error)
	ModifyLedgerName(ws *workspace.Workspace, ledgerId, ledgerName string) error
	ListAllLedger(ws *workspace.Workspace) ([]models.Ledger, error)
	QueryLedgerById(ws *workspace.Workspace, ledgerId string) (*models.Ledger, error)
	DeleteLedgerById(ws *workspace.Workspace, ledgerId string) error
}

var _ LedgerService = &ledgerServiceImpl{}

type ledgerServiceImpl struct {
	ledgerDao dao.LedgerDao
	trDao     dao.TransactionRecordDao
	trTagDao  dao.TrTagDao
}

// CreateLedger 创建成功返回创建账本id
func (l *ledgerServiceImpl) CreateLedger(ws *workspace.Workspace, ledgerName string) (string, error) {
	log := ws.GetLogger()
	log.Infof("start to create ledger, name: %s", ledgerName)
	ledger := &models.Ledger{
		ID:   util.GetUUID(),
		Name: ledgerName,
	}

	if err := l.ledgerDao.CreateLedger(ws, ledger); err != nil {
		log.Errorf("create ledger failed, name: %s, err: %v", ledgerName, err)
		return "", err
	}

	log.Infof("create ledger success, name: %s", ledgerName)
	return ledger.ID, nil
}

// ModifyLedgerName 修改指定账本的名称
func (l *ledgerServiceImpl) ModifyLedgerName(ws *workspace.Workspace, ledgerId, ledgerName string) error {
	log := ws.GetLogger()
	log.Infof("start to modify ledger name, id: %s, new name: %s", ledgerId, ledgerName)

	ledger := &models.Ledger{
		ID:   ledgerId,
		Name: ledgerName,
	}

	if err := l.ledgerDao.ModifyLedgerName(ws, ledger); err != nil {
		log.Errorf("modify ledger name failed, id: %s, err: %v", ledgerId, err)
		return err
	}

	log.Infof("modify ledger name success")
	return nil
}

// ListAllLedger 查询所有账本
func (l *ledgerServiceImpl) ListAllLedger(ws *workspace.Workspace) ([]models.Ledger, error) {
	log := ws.GetLogger()
	log.Infof("start to list all ledgers")

	ledgers, err := l.ledgerDao.ListAllLedger(ws)
	if err != nil {
		log.Errorf("list all ledgers failed, err: %v", err)
		return nil, err
	}

	log.Infof("end to list all ledgers, len: %d", len(ledgers))
	return ledgers, nil
}

// QueryLedgerById 查询单个账本
func (l *ledgerServiceImpl) QueryLedgerById(ws *workspace.Workspace, ledgerId string) (*models.Ledger, error) {
	log := ws.GetLogger()
	log.Infof("start to query ledger by id, id: %s", ledgerId)

	ledger, err := l.ledgerDao.QueryLedgerById(ws, ledgerId)
	if err != nil {
		log.Errorf("query ledger by id failed, id: %s, err: %v", ledgerId, err)
		return nil, err
	}

	log.Infof("end to query ledger by id, id: %s", ledgerId)
	return ledger, nil
}

// DeleteLedgerById deletes a ledger and all its associated transaction records and tags in a transaction.
func (l *ledgerServiceImpl) DeleteLedgerById(ws *workspace.Workspace, ledgerId string) error {
	log := ws.GetLogger()
	log.Infof("start to delete ledger by id, id: %s", ledgerId)

	err := ws.Transaction(func(tx *workspace.Workspace) error {
		// Delete all tags for this ledger's transactions
		if err := l.trTagDao.DeleteTrTagByLedgerId(tx, ledgerId); err != nil {
			return fmt.Errorf("delete tr tags: %w", err)
		}

		// Count and delete all transaction records
		cnt, err := l.trDao.CountTrByLedgerId(tx, ledgerId)
		if err != nil {
			return fmt.Errorf("count trs: %w", err)
		}
		log.Infof("will delete trs by ledger id: %s, count: %d", ledgerId, cnt)

		if err := l.trDao.DeleteAllTrByLedgerId(tx, ledgerId); err != nil {
			return fmt.Errorf("delete trs: %w", err)
		}

		// Delete the ledger
		if err := l.ledgerDao.DeleteLedgerById(tx, ledgerId); err != nil {
			return fmt.Errorf("delete ledger: %w", err)
		}
		return nil
	})

	if err != nil {
		log.Errorf("delete ledger by id failed, id: %s, err: %v", ledgerId, err)
		return err
	}

	log.Infof("end to delete ledger by id, id: %s", ledgerId)
	return nil
}
