package service

import (
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
	ws.GetLogger().Infof("start to create ledger, name: %s", ledgerName)
	// 账本名称可以重复，不需要校验账本名称
	ledger := &models.Ledger{
		ID:   util.GetUUID(),
		Name: ledgerName,
	}

	if err := l.ledgerDao.CreateLedger(ws, ledger); err != nil {
		ws.GetLogger().Errorf("create ledger failed, name: %s, err: %v", ledgerName, err)
		return "", err
	}

	ws.GetLogger().Infof("create ledger success, name: %s", ledgerName)
	return ledger.ID, nil
}

// ModifyLedgerName 修改指定账本的名称
func (l *ledgerServiceImpl) ModifyLedgerName(ws *workspace.Workspace, ledgerId, ledgerName string) error {
	ws.GetLogger().Infof("start to modify ledger name, id: %s, new name: %s", ledgerId, ledgerName)

	ledger := &models.Ledger{
		ID:   ledgerId,
		Name: ledgerName,
	}

	if err := l.ledgerDao.ModifyLedgerName(ws, ledger); err != nil {
		ws.GetLogger().Errorf("modify ledger name failed, id: %s, err: %v", ledgerId, err)
		return err
	}

	ws.GetLogger().Infof("modify ledger name success")
	return nil
}

// ListAllLedger 查询所有账本
func (l *ledgerServiceImpl) ListAllLedger(ws *workspace.Workspace) ([]models.Ledger, error) {
	ws.GetLogger().Infof("start to list all ledgers")

	ledgers, err := l.ledgerDao.ListAllLedger(ws)
	if err != nil {
		ws.GetLogger().Errorf("list all ledgers failed, err: %v", err)
		return nil, err
	}

	ws.GetLogger().Infof("end to list all ledgers, len: %d", len(ledgers))
	return ledgers, nil
}

// QueryLedgerById 查询单个账本
func (l *ledgerServiceImpl) QueryLedgerById(ws *workspace.Workspace, ledgerId string) (*models.Ledger, error) {
	ws.GetLogger().Infof("start to query ledger by id, id: %s", ledgerId)

	ledger, err := l.ledgerDao.QueryLedgerById(ws, ledgerId)
	if err != nil {
		ws.GetLogger().Errorf("query ledger by id failed, id: %s, err: %v", ledgerId, err)
		return nil, err
	}

	ws.GetLogger().Infof("end to query ledger by id, id: %s", ledgerId)
	return ledger, nil
}

func (l *ledgerServiceImpl) DeleteLedgerById(ws *workspace.Workspace, ledgerId string) error {
	ws.GetLogger().Infof("start to delete ledger by id, id: %s", ledgerId)

	// 删除账本中消费记录的所有tag
	if err := l.trTagDao.DeleteTrTagByLedgerId(ws, ledgerId); err != nil {
		ws.GetLogger().Errorf("delete all trTags by id failed, id: %s, err: %v", ledgerId, err)
		return err
	}

	// 删除账本中的所有消费记录
	cnt, err := l.trDao.CountTrByLedgerId(ws, ledgerId)
	if err != nil {
		ws.GetLogger().Errorf("get count of trs from ledger by id failed, id: %s, err: %v", ledgerId, err)
		return err
	}
	ws.GetLogger().Infof("will delete trs by ledger id: %s, count: %d", ledgerId, cnt)

	if err := l.trDao.DeleteAllTrByLedgerId(ws, ledgerId); err != nil {
		ws.GetLogger().Errorf("delete all tr by ledger id failed, id: %s, err: %v", ledgerId, err)
		return err
	}

	// 删除账本
	if err := l.ledgerDao.DeleteLedgerById(ws, ledgerId); err != nil {
		ws.GetLogger().Errorf("delete ledger by id failed, id: %s, err: %v", ledgerId, err)
		return err
	}

	ws.GetLogger().Infof("end to delete ledger by id, id: %s", ledgerId)
	return nil
}
