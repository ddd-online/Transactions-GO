package service

import (
	"sync"

	"github.com/billadm/dao"
	"github.com/billadm/models"
	"github.com/billadm/workspace"
)

var (
	categoryService     CategoryService
	categoryServiceOnce sync.Once
)

func GetCategoryService() CategoryService {
	if categoryService != nil {
		return categoryService
	}

	categoryServiceOnce.Do(func() {
		categoryService = &categoryServiceImpl{
			categoryDao: dao.GetCategoryDao(),
		}
	})

	return categoryService
}

type CategoryService interface {
	QueryCategory(ws *workspace.Workspace, trType string) ([]models.Category, error)
}

var _ CategoryService = &categoryServiceImpl{}

type categoryServiceImpl struct {
	categoryDao dao.CategoryDao
}

func (c *categoryServiceImpl) QueryCategory(ws *workspace.Workspace, trType string) ([]models.Category, error) {
	ws.GetLogger().Infof("start to query category by %s", trType)
	categories, err := c.categoryDao.QueryCategory(ws, trType)
	if err != nil {
		return nil, err
	}

	ws.GetLogger().Infof("query category success, length: %v", categories)
	return categories, nil
}
