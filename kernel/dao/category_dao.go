package dao

import (
	"sync"

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
}

var _ CategoryDao = &categoryDaoImpl{}

type categoryDaoImpl struct{}

func (c *categoryDaoImpl) QueryCategory(ws *workspace.Workspace, trType string) ([]models.Category, error) {
	categories := make([]models.Category, 0)
	db := ws.GetDb()
	if trType != "all" {
		db = db.Where("transaction_type = ?", trType)
	}

	if err := db.Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}
