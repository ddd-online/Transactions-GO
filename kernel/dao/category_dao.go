package dao

import (
	"sync"

	"github.com/billadm/constant"
	"github.com/billadm/models"
	"github.com/billadm/workspace"
)

var (
	categoryDao     CategoryDao
	categoryDaoOnce sync.Once
)

func GetCategoryDao() CategoryDao {
	if categoryDao != nil {
		return categoryDao
	}
	categoryDaoOnce.Do(func() {
		categoryDao = &categoryDaoImpl{}
	})
	return categoryDao
}

type CategoryDao interface {
	QueryCategory(ws *workspace.Workspace, trType string) ([]models.Category, error)
	CreateCategory(ws *workspace.Workspace, category *models.Category) error
	DeleteCategory(ws *workspace.Workspace, name string, transactionType string) error
	IsCategoryInUse(ws *workspace.Workspace, ledgerId string, category string) (bool, error)
	UpdateCategorySort(ws *workspace.Workspace, name string, transactionType string, sortOrder int) error
	GetMaxSortOrder(ws *workspace.Workspace, transactionType string) (int, error)
	CountRecordsByCategory(ws *workspace.Workspace, ledgerId string, category string) (int64, error)
}

var _ CategoryDao = &categoryDaoImpl{}

type categoryDaoImpl struct{}

func (c *categoryDaoImpl) QueryCategory(ws *workspace.Workspace, trType string) ([]models.Category, error) {
	categories := make([]models.Category, 0)
	db := ws.GetDb()
	if trType != constant.All {
		db = db.Where("transaction_type = ?", trType)
	}

	if err := db.Order("sort_order ASC, name DESC").Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (c *categoryDaoImpl) CreateCategory(ws *workspace.Workspace, category *models.Category) error {
	if err := ws.GetDb().Create(category).Error; err != nil {
		return err
	}
	return nil
}

func (c *categoryDaoImpl) DeleteCategory(ws *workspace.Workspace, name string, transactionType string) error {
	if err := ws.GetDb().
		Where("name = ? AND transaction_type = ?", name, transactionType).
		Delete(&models.Category{}).Error; err != nil {
		return err
	}
	return nil
}

func (c *categoryDaoImpl) IsCategoryInUse(ws *workspace.Workspace, ledgerId string, category string) (bool, error) {
	var count int64
	err := ws.GetDb().Model(&models.TransactionRecord{}).
		Where("ledger_id = ? AND category = ?", ledgerId, category).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (c *categoryDaoImpl) UpdateCategorySort(ws *workspace.Workspace, name string, transactionType string, sortOrder int) error {
	if err := ws.GetDb().
		Model(&models.Category{}).
		Where("name = ? AND transaction_type = ?", name, transactionType).
		Update("sort_order", sortOrder).Error; err != nil {
		return err
	}
	return nil
}

func (c *categoryDaoImpl) GetMaxSortOrder(ws *workspace.Workspace, transactionType string) (int, error) {
	var maxSortOrder int
	err := ws.GetDb().Model(&models.Category{}).
		Where("transaction_type = ?", transactionType).
		Select("COALESCE(MAX(sort_order), 0)").
		Scan(&maxSortOrder).Error
	if err != nil {
		return 0, err
	}
	return maxSortOrder, nil
}

func (c *categoryDaoImpl) CountRecordsByCategory(ws *workspace.Workspace, ledgerId string, category string) (int64, error) {
	var count int64
	err := ws.GetDb().Model(&models.TransactionRecord{}).
		Where("ledger_id = ? AND category = ?", ledgerId, category).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}