package dao

import (
	"sync"

	"github.com/billadm/models"
	"github.com/billadm/workspace"
)

var (
	chartDao     ChartDao
	chartDaoOnce sync.Once
)

func GetChartDao() ChartDao {
	if chartDao != nil {
		return chartDao
	}

	chartDaoOnce.Do(func() {
		chartDao = &chartDaoImpl{}
	})

	return chartDao
}

type ChartDao interface {
	Create(ws *workspace.Workspace, chart *models.Chart) error
	DeleteById(ws *workspace.Workspace, chartId string) error
	ListByLedgerId(ws *workspace.Workspace, ledgerId string) ([]*models.Chart, error)
	GetById(ws *workspace.Workspace, chartId string) (*models.Chart, error)
	Update(ws *workspace.Workspace, chart *models.Chart) error
	GetMaxSortOrder(ws *workspace.Workspace, ledgerId string) (int, error)
}

var _ ChartDao = &chartDaoImpl{}

type chartDaoImpl struct{}

func (t *chartDaoImpl) Create(ws *workspace.Workspace, chart *models.Chart) error {
	if err := ws.GetDb().Create(chart).Error; err != nil {
		return err
	}
	return nil
}

func (t *chartDaoImpl) DeleteById(ws *workspace.Workspace, chartId string) error {
	if err := ws.GetDb().
		Where("chart_id = ?", chartId).
		Delete(&models.Chart{}).Error; err != nil {
		return err
	}
	return nil
}

func (t *chartDaoImpl) ListByLedgerId(ws *workspace.Workspace, ledgerId string) ([]*models.Chart, error) {
	charts := make([]*models.Chart, 0)
	if err := ws.GetDb().
		Where("ledger_id = ?", ledgerId).
		Order("is_preset DESC, sort_order ASC, created_at DESC").
		Find(&charts).Error; err != nil {
		return nil, err
	}
	return charts, nil
}

func (t *chartDaoImpl) GetById(ws *workspace.Workspace, chartId string) (*models.Chart, error) {
	var chart models.Chart
	if err := ws.GetDb().
		Where("chart_id = ?", chartId).
		First(&chart).Error; err != nil {
		return nil, err
	}
	return &chart, nil
}

func (t *chartDaoImpl) Update(ws *workspace.Workspace, chart *models.Chart) error {
	if err := ws.GetDb().Save(chart).Error; err != nil {
		return err
	}
	return nil
}

func (t *chartDaoImpl) GetMaxSortOrder(ws *workspace.Workspace, ledgerId string) (int, error) {
	var maxSortOrder int
	err := ws.GetDb().Model(&models.Chart{}).
		Where("ledger_id = ?", ledgerId).
		Select("COALESCE(MAX(sort_order), 0)").
		Scan(&maxSortOrder).Error
	if err != nil {
		return 0, err
	}
	return maxSortOrder, nil
}
